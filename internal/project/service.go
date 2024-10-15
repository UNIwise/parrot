//go:generate mockgen --source=service.go -destination=service_mock.go -package=project

package project

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	gosundheit "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	"github.com/alitto/pond"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/uniwise/parrot/internal/cache"
	"github.com/uniwise/parrot/internal/storage"
	"github.com/uniwise/parrot/pkg/poedit"
	"golang.org/x/sync/semaphore"
)

var ErrVersionAlreadyExist = errors.New("Project version already exists")

type Translation struct {
	TTL      time.Duration
	Checksum string
	Data     []byte
}

type Service interface {
	GetTranslation(ctx context.Context, projectID int, languageCode, format, version string) (trans *Translation, err error)
	PurgeTranslation(ctx context.Context, projectID int, languageCode string) (err error)
	PurgeProject(ctx context.Context, projectID int) (err error)
	RegisterChecks(h gosundheit.Health) (err error)
	GetAllProjects(ctx context.Context) ([]Project, error)
	GetProjectByID(ctx context.Context, id int) (*Project, error)
	GetProjectVersions(ctx context.Context, projectID int) ([]Version, error)
	DeleteProjectVersionByIDAndProjectID(ctx context.Context, ID string, projectID uint) error
	CreateLanguagesVersion(ctx context.Context, projectID int, name string) error
}

type ServiceImpl struct {
	Logger            *logrus.Entry
	Client            poedit.Client
	Cache             cache.Cache
	RenewalThreshold  time.Duration
	PreFetchSemaphore *semaphore.Weighted
	storage           storage.Storage
	generateUUID      func() (string, error)
	generateTimestamp func() int64
	getContentMetaMap func() map[string]poedit.ContentMeta
}

func NewService(cli poedit.Client, storage storage.Storage, cache cache.Cache, renewalThreshold time.Duration, entry *logrus.Entry) *ServiceImpl {
	return &ServiceImpl{
		Logger:            entry,
		Client:            cli,
		Cache:             cache,
		RenewalThreshold:  renewalThreshold,
		PreFetchSemaphore: semaphore.NewWeighted(1),
		storage:           storage,
		generateUUID:      GenerateUUID,
		generateTimestamp: GenerateTimestamp,
		getContentMetaMap: GetContentMetaMap,
	}
}

func (s *ServiceImpl) GetTranslation(ctx context.Context, projectID int, languageCode, format, version string) (*Translation, error) {
	item, err := s.Cache.GetTranslation(ctx, projectID, languageCode, format, version)
	if err != nil && !errors.Is(err, cache.ErrCacheMiss) {
		return nil, err
	}
	if err == nil {
		expiresAt := item.CreatedAt.Add(s.Cache.GetTTL())

		if time.Until(expiresAt) < s.RenewalThreshold {
			go func() {
				if !s.PreFetchSemaphore.TryAcquire(1) {
					return
				}

				defer s.PreFetchSemaphore.Release(1)

				s.Logger.Debugf("Pre-fetching language %s format %s for project %d", languageCode, format, projectID)

				_, _, err := s.fetchAndCacheTranslation(context.Background(), projectID, languageCode, format, version)
				if err != nil {
					s.Logger.Errorf("Failed to pre-fetch language %s format %s for project %d", languageCode, format, projectID)
				}
			}()
		}

		return &Translation{
			TTL:      s.Cache.GetTTL(),
			Checksum: item.Checksum,
			Data:     item.Data,
		}, nil
	}

	data, checksum, err := s.fetchAndCacheTranslation(ctx, projectID, languageCode, format, version)
	if err != nil {
		return nil, err
	}

	return &Translation{
		TTL:      s.Cache.GetTTL(),
		Checksum: checksum,
		Data:     data,
	}, nil
}

func (s *ServiceImpl) PurgeTranslation(ctx context.Context, projectID int, languageCode string) error {
	return errors.New("Not implemented")
}

func (s *ServiceImpl) PurgeProject(ctx context.Context, projectID int) error {
	return errors.New("Not implemented")
}

func (s *ServiceImpl) fetchAndCacheTranslation(ctx context.Context, projectID int, languageCode, format, version string) ([]byte, string, error) {
	data, err := s.fetchTranslation(ctx, projectID, languageCode, format, version)
	if err != nil {
		return nil, "", err
	}

	checksum, err := s.Cache.SetTranslation(ctx, projectID, languageCode, format, version, data)
	if err != nil {
		return nil, "", err
	}

	return data, checksum, nil
}

func (s *ServiceImpl) fetchTranslation(ctx context.Context, projectID int, languageCode, format, version string) ([]byte, error) {
	if version == "latest" {
		return s.getProjectTranslationFromPOedit(ctx, projectID, languageCode, format)
	}
	return s.getProjectTranslationFromS3(ctx, projectID, version, languageCode, format)
}

func (s *ServiceImpl) getProjectTranslationFromS3(ctx context.Context, projectID int, versionName, languageCode, format string) ([]byte, error) {
	s3Output, err := s.storage.ListObjects(ctx, fmt.Sprintf("%d/", projectID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project versions")
	}

	s3Key := ""
	for _, object := range s3Output.CommonPrefixes {
		// Example Prefix : {projectID}/{versionID_versionName_timestamp}/
		// 720964/61ded6dc-c8b7-4d4e-aa70-cd37dd1216b3_v2_123456789/
		prefixData := strings.Split(aws.ToString(object.Prefix), "/")

		versionData := strings.Split(prefixData[1], "_")
		version := versionData[1]

		if version == versionName {
			contentMeta, _ := poedit.GetContentMeta(format)
			s3Key = fmt.Sprintf("%s%s.%s", aws.ToString(object.Prefix), languageCode, contentMeta.Extension)
			break
		}
	}

	if s3Key != "" {

		reader, err := s.storage.GetObject(ctx, s3Key)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get object from S3")
		}

		data, err := io.ReadAll(reader.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to read object from S3")
		}

		return data, nil
	}

	return nil, &ErrLanguageNotFoundInStorage{
		ProjectID:    projectID,
		LanguageCode: languageCode,
		Version:      versionName,
	}
}

func (s *ServiceImpl) getProjectTranslationFromPOedit(ctx context.Context, projectID int, languageCode, format string) ([]byte, error) {
	resp, err := s.Client.ExportProject(ctx, poedit.ExportProjectRequest{
		ID:       projectID,
		Language: languageCode,
		Type:     format,
		Filters:  []string{"translated"},
	})
	if err != nil {
		return nil, err
	}

	// TODO: Make use of injected http client, to support timeouts
	d, err := http.Get(resp.Result.URL)
	if err != nil {
		return nil, err
	}
	defer d.Body.Close()

	if d.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Response code '%d' from download GET", d.StatusCode)
	}

	data, err := io.ReadAll(d.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *ServiceImpl) RegisterChecks(h gosundheit.Health) error {
	c, err := checks.NewPingCheck("cache", s.Cache)
	if err != nil {
		return errors.Wrap(err, "Failed to instantiate cache healthcheck")
	}

	if err := h.RegisterCheck(c, gosundheit.ExecutionPeriod(time.Second*10)); err != nil {
		return errors.Wrap(err, "Failed to register cache healthcheck")
	}

	return nil
}

func (s *ServiceImpl) GetAllProjects(ctx context.Context) ([]Project, error) {
	var projects []Project

	projectsResponse, err := s.Client.ListProjects(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get projects from poeditor")
	}

	for _, project := range projectsResponse.Result.Projects {
		s3Output, err := s.storage.ListObjects(ctx, fmt.Sprintf("%d/", project.ID))
		if err != nil {
			return nil, errors.Wrap(err, "failed to get project versions from S3")
		}

		numberOfVersions := len(s3Output.CommonPrefixes)

		projects = append(projects, Project{
			ID:               project.ID,
			Name:             project.Name,
			NumberOfVersions: numberOfVersions,
			CreatedAt:        project.Created,
		})
	}

	return projects, nil
}

func (s *ServiceImpl) GetProjectByID(ctx context.Context, id int) (*Project, error) {
	projectResponse, err := s.Client.ViewProject(ctx, poedit.ViewProjectRequest{
		ID: id,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to view project")
	}

	projectResult := projectResponse.Result.Project
	s3Output, err := s.storage.ListObjects(ctx, fmt.Sprintf("%d/", id))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project versions")
	}

	numberOfVersions := len(s3Output.CommonPrefixes)

	project := &Project{
		ID:               projectResult.ID,
		Name:             projectResult.Name,
		NumberOfVersions: numberOfVersions,
		CreatedAt:        projectResult.Created,
	}

	return project, nil
}

func (s *ServiceImpl) GetProjectVersions(ctx context.Context, projectID int) ([]Version, error) {
	var versions []Version

	s3Output, err := s.storage.ListObjects(ctx, fmt.Sprintf("%d/", projectID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get project versions")
	}

	for _, object := range s3Output.CommonPrefixes {
		// Example Prefix : {projectID}/{versionID_versionName_timestamp}/
		// 720964/61ded6dc-c8b7-4d4e-aa70-cd37dd1216b3_v2_123456789/
		prefixData := strings.Split(aws.ToString(object.Prefix), "/")

		versionData := strings.Split(prefixData[1], "_")
		versionID := versionData[0]
		versionName := versionData[1]
		versionTimestamp, err := strconv.ParseInt(versionData[2], 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse timestamp")
		}

		versions = append(versions, Version{
			ID:        versionID,
			Name:      versionName,
			CreatedAt: time.Unix(versionTimestamp, 0),
		})
	}

	return versions, nil
}

func (s *ServiceImpl) DeleteProjectVersionByIDAndProjectID(ctx context.Context, ID string, projectID uint) error {
	storageKey := fmt.Sprintf("%d/%s", projectID, ID)
	if err := s.storage.DeleteObjects(ctx, storageKey); err != nil {
		return errors.Wrap(err, "Failed to delete project version in S3")
	}

	return nil
}

func (s *ServiceImpl) checkProjectVersionExistInS3(ctx context.Context, projectID int, name string) (bool, error) {
	s3Output, err := s.storage.ListObjects(ctx, fmt.Sprintf("%d/", projectID))
	if err != nil {
		return false, errors.Wrap(err, "failed to get project versions")
	}

	for _, object := range s3Output.CommonPrefixes {
		// Example Prefix : {projectID}/{versionID_versionName_timestamp}/
		// 720964/61ded6dc-c8b7-4d4e-aa70-cd37dd1216b3_1.0.0_123456789/
		prefixData := strings.Split(aws.ToString(object.Prefix), "/")

		versionData := strings.Split(prefixData[1], "_")
		versionName := versionData[1]

		if versionName == name {
			return true, nil
		}
	}

	return false, nil
}

func (s *ServiceImpl) CreateLanguagesVersion(ctx context.Context, projectID int, name string) error {
	exists, err := s.checkProjectVersionExistInS3(ctx, projectID, name)
	if err != nil {
		return errors.Wrap(err, "Failed to check project version existence in S3")
	}
	if exists {
		return ErrVersionAlreadyExist
	}
	languagesResponse, err := s.Client.ListProjectLanguages(ctx, poedit.ListProjectLanguagesRequest{
		ID: projectID,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to list project languages")
	}

	uuid, err := s.generateUUID()
	if err != nil {
		return errors.Wrap(err, "Failed to generate UUID")
	}

	timeStamp := s.generateTimestamp()
	contentMetaMap := s.getContentMetaMap()

	type TranslationObject struct {
		S3Key       string
		Reader      *bytes.Reader
		Meta        map[string]string
		ContentType string
	}

	jobErrors := []error{}
	pool := pond.New(10, 1000) // 10 workers, 1000 buffered jobs

	versionKey := fmt.Sprintf("%d/%s_%s_%d", projectID, uuid, name, timeStamp)

	for _, language := range languagesResponse.Result.Languages {
		for contentMetaKey, contentMeta := range contentMetaMap {
			pool.Submit(func() {
				resp, err := s.Client.ExportProject(ctx, poedit.ExportProjectRequest{
					ID:       projectID,
					Language: language.Code,
					Type:     contentMetaKey,
					Filters:  []string{"translated"},
				})
				if err != nil {
					jobErrors = append(jobErrors, errors.Wrap(err, "Failed to export project"))
				}

				d, err := http.Get(resp.Result.URL)
				if err != nil {
					jobErrors = append(jobErrors, errors.Wrap(err, "Failed to download project language file"))
				}
				defer d.Body.Close()

				if d.StatusCode != http.StatusOK {
					jobErrors = append(jobErrors, errors.Errorf("Response code '%d' from download GET", d.StatusCode))
				}

				s3Key := fmt.Sprintf("%s/%s.%s", versionKey, language.Code, contentMeta.Extension)

				meta := map[string]string{
					"project":     fmt.Sprintf("%d", projectID),
					"lang":        language.Code,
					"versionName": name,
					"format":      contentMeta.Extension,
				}

				l := s.Logger.WithFields(logrus.Fields{
					"project":     projectID,
					"lang":        language.Code,
					"versionName": name,
					"format":      contentMeta.Extension,
					"key":         s3Key,
				})

				l.Debug("Uploading")

				if err := s.storage.PutObject(ctx, s3Key, d.Body, meta, contentMeta.Type); err != nil {
					jobErrors = append(jobErrors, errors.Wrap(err, "Failed to upload project language file to S3"))
				}

				l.Debug("Done uploading")
			})
		}
	}

	pool.StopAndWait()

	if len(jobErrors) > 0 {
		if err := s.storage.DeleteObjects(ctx, versionKey); err != nil {
			return errors.Wrap(err, "Failed to delete project version in S3")
		}

		return errors.New("Failed to create language versions")
	}

	return nil
}
