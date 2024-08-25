package project

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/uniwise/parrot/pkg/connectors/database"
)

var (
	queryGetAllProjects                    string = regexp.QuoteMeta("SELECT projects.id, projects.name, COUNT(versions.id) as number_of_versions, projects.created_at FROM `projects` LEFT JOIN versions ON projects.id = versions.project_id GROUP BY `projects`.`id`")
	queryGetProjectById                    string = regexp.QuoteMeta("SELECT projects.id, projects.name, COUNT(versions.id) as number_of_versions, projects.created_at FROM `projects` LEFT JOIN versions ON projects.id = versions.project_id WHERE `projects`.`id` = ? GROUP BY `projects`.`id` ORDER BY `projects`.`id` LIMIT ?")
	queryGetProjectVersions                string = regexp.QuoteMeta("SELECT * FROM `versions` WHERE project_id = ?")
	queryGetProjectVersionByIdAndProjectID string = regexp.QuoteMeta("SELECT * FROM `versions` WHERE project_id = ? AND id = ? ORDER BY `versions`.`id` LIMIT ?")
	queryDeleteProjectVersionTransaction   string = regexp.QuoteMeta("DELETE FROM `versions` WHERE `versions`.`id` = ?")
)

func TestRepositoryGetAllProjects(t *testing.T) {
	t.Parallel()

	db, sql := database.NewMockClient(t)
	repository := NewRepository(db)

	timestamp := time.Now()

	t.Run("GetAllProjects, success", func(t *testing.T) {
		sql.ExpectQuery(queryGetAllProjects).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "number_of_versions", "created_at"}).
				AddRow(testID, "testName", testNumberOfVersions, timestamp),
		)

		projects, err := repository.GetAllProjects(context.Background())

		assert.NoError(t, err)
		assert.Len(t, projects, 1)
		assert.Equal(t, testID, projects[0].ID)
		assert.Equal(t, "testName", projects[0].Name)
		assert.Equal(t, testNumberOfVersions, projects[0].NumberOfVersions)
		assert.Equal(t, timestamp, projects[0].CreatedAt)
		assert.NoError(t, sql.ExpectationsWereMet())
	})

	t.Run("GetAllProjects, error", func(t *testing.T) {
		sql.ExpectQuery(queryGetAllProjects).WillReturnError(assert.AnError)

		projects, err := repository.GetAllProjects(context.Background())

		assert.Error(t, err)
		assert.Nil(t, projects)
		assert.NoError(t, sql.ExpectationsWereMet())
	})
}

func TestRepositoryGetProjectById(t *testing.T) {
	t.Parallel()

	db, sql := database.NewMockClient(t)
	repository := NewRepository(db)

	timestamp := time.Now()

	t.Run("GetProject, success", func(t *testing.T) {
		sql.ExpectQuery(queryGetProjectById).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "number_of_versions", "created_at"}).
				AddRow(testID, "testName", testNumberOfVersions, timestamp),
		)

		project, err := repository.GetProjectByID(context.Background(), int(testID))

		assert.NoError(t, err)
		assert.Equal(t, testID, (*project).ID)
		assert.Equal(t, "testName", (*project).Name)
		assert.Equal(t, testNumberOfVersions, (*project).NumberOfVersions)
		assert.Equal(t, timestamp, (*project).CreatedAt)
		assert.NoError(t, sql.ExpectationsWereMet())
	})

	t.Run("GetProject, error", func(t *testing.T) {
		sql.ExpectQuery(queryGetProjectById).WillReturnError(assert.AnError)

		project, err := repository.GetProjectByID(context.Background(), int(testID))

		assert.Error(t, err)
		assert.Nil(t, project)
		assert.NoError(t, sql.ExpectationsWereMet())
	})
}

func TestRepositoryGetProjectVersions(t *testing.T) {
	t.Parallel()

	db, sql := database.NewMockClient(t)
	repository := NewRepository(db)

	timestamp := time.Now()

	t.Run("GetProjectVersions, success", func(t *testing.T) {
		sql.ExpectQuery(queryGetProjectVersions).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "project_id", "created_at"}).
				AddRow(testID, "testName", testProjectID, timestamp),
		)

		versions, err := repository.GetProjectVersions(context.Background(), int(testProjectID))

		assert.NoError(t, err)
		assert.Len(t, versions, 1)
		assert.Equal(t, testID, versions[0].ID)
		assert.Equal(t, "testName", versions[0].Name)
		assert.Equal(t, testProjectID, versions[0].ProjectID)
		assert.Equal(t, timestamp, versions[0].CreatedAt)
		assert.NoError(t, sql.ExpectationsWereMet())
	})

	t.Run("GetProjectVersions, error", func(t *testing.T) {
		sql.ExpectQuery(queryGetProjectVersions).WillReturnError(assert.AnError)

		versions, err := repository.GetProjectVersions(context.Background(), int(testProjectID))

		assert.Error(t, err)
		assert.Nil(t, versions)
		assert.NoError(t, sql.ExpectationsWereMet())
	})
}

func TestRepositoryGetProjectVersionByIDAndProjectID(t *testing.T) {
	t.Parallel()

	db, sql := database.NewMockClient(t)
	repository := NewRepository(db)

	timestamp := time.Now()

	t.Run("GetProjectVersionByIDAndProjectID, success", func(t *testing.T) {
		sql.ExpectQuery(queryGetProjectVersionByIdAndProjectID).WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "project_id", "storage_key", "created_at"}).
				AddRow(testID, "testName", testProjectID, testStorageKey, timestamp),
		)

		version, err := repository.GetVersionByIDAndProjectID(context.Background(), testID, testProjectID)

		assert.NoError(t, err)
		assert.Equal(t, testID, version.ID)
		assert.Equal(t, "testName", version.Name)
		assert.Equal(t, testProjectID, version.ProjectID)
		assert.Equal(t, testStorageKey, version.StorageKey)
		assert.Equal(t, timestamp, version.CreatedAt)
		assert.NoError(t, sql.ExpectationsWereMet())
	})

	t.Run("GetProjectVersionByIDAndProjectID, error", func(t *testing.T) {
		sql.ExpectQuery(queryGetProjectVersionByIdAndProjectID).WillReturnError(assert.AnError)

		version, err := repository.GetVersionByIDAndProjectID(context.Background(), testID, testProjectID)

		assert.Error(t, err)
		assert.Nil(t, version)
		assert.NoError(t, sql.ExpectationsWereMet())
	})
}

func TestRepositoryDeleteVersionByIDTransaction(t *testing.T) {
	t.Parallel()

	db, sql := database.NewMockClient(t)
	repository := NewRepository(db)

	t.Run("DeleteVersionByIDTransaction, fail, noneDeleted", func(t *testing.T) {
		sql.ExpectBegin()
		sql.ExpectExec(queryDeleteProjectVersionTransaction).
			WillReturnResult(sqlmock.NewResult(0, 0))
		sql.ExpectRollback()

		tx, err := repository.DeleteVersionByIDTransaction(testCtx, testID)

		assert.ErrorIs(t, err, ErrNotDeleted)
		assert.Nil(t, tx)
		assert.NoError(t, sql.ExpectationsWereMet())
	})

	t.Run("DeleteVersionByIDTransaction, error", func(t *testing.T) {
		sql.ExpectBegin()
		sql.ExpectExec(queryDeleteProjectVersionTransaction).
			WillReturnError(errTest)
		sql.ExpectRollback()

		tx, err := repository.DeleteVersionByIDTransaction(testCtx, testID)

		assert.ErrorIs(t, err, errTest)
		assert.Nil(t, tx)
		assert.NoError(t, sql.ExpectationsWereMet())
	})

	t.Run("DeleteVersionByIDTransaction, success", func(t *testing.T) {
		sql.ExpectBegin()
		sql.ExpectExec(queryDeleteProjectVersionTransaction).
			WillReturnResult(sqlmock.NewResult(1, 1))

		tx, err := repository.DeleteVersionByIDTransaction(testCtx, testID)

		assert.NoError(t, err)
		assert.NotNil(t, tx)
		assert.NoError(t, sql.ExpectationsWereMet())
	})
}
