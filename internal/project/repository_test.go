package project

import (
	"regexp"
	"testing"
	"time"
	"context"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/uniwise/parrot/pkg/connectors/database"
)

var (
	queryGetAllProjects string = regexp.QuoteMeta("SELECT projects.id, projects.name, COUNT(versions.id) as number_of_versions, projects.created_at FROM `projects` LEFT JOIN versions ON projects.id = versions.project_id GROUP BY `projects`.`id`")
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
		assert.Len(t, *projects, 1)
		assert.Equal(t, testID, (*projects)[0].ID)
		assert.Equal(t, "testName", (*projects)[0].Name)
		assert.Equal(t, testNumberOfVersions, (*projects)[0].NumberOfVersions)
		assert.Equal(t, timestamp, (*projects)[0].CreatedAt)
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
