package database

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)

	return ok
}

func NewMockClient(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mockSQL, err := sqlmock.New()
	assert.NoError(t, err)

	columns := []string{"version"}
	mockSQL.ExpectQuery("SELECT VERSION()").WithArgs().WillReturnRows(
		mockSQL.NewRows(columns).FromCSVString("1"),
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return db, mockSQL
}