package kumload

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func TestNewMigration(t *testing.T) {
	conn, _, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	db, err := gorm.Open("postgres", conn)
	if err != nil {
		t.Error(err)
	}
	cfg := &Config{}
	expected := &migrationImpl{
		db:     db,
		config: cfg,
	}

	result := NewMigration(db, cfg)
	if result == nil {
		t.Error("unexpected nil result")
	}

	if m, ok := result.(*migrationImpl); ok {
		if m.config != expected.config {
			t.Errorf("expected: %+v, got: %+v", expected.config, m.config)
		}
		if m.db != expected.db {
			t.Errorf("expected: %+v, got: %+v", expected.db, m.db)
		}
	}
}

func TestMigration_Run(t *testing.T) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}
	mock.MatchExpectationsInOrder(false)

	db, err := gorm.Open("postgres", conn)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	count := int64(10)
	migration := &migrationImpl{
		db: db,
		config: &Config{
			Batch: Batch{
				Size:     2,
				Interval: 0,
			},
			Database: Database{
				Source: DatabaseDetail{
					Name:       "user_service",
					Table:      "sessions",
					PrimaryKey: "id",
					OrderKey:   "created_at",
				},
				Target: DatabaseDetail{
					Name:       "authorization_service",
					Table:      "sessions",
					PrimaryKey: "id",
				},
			},
			Mappings: map[string]string{
				"username": "username",
			},
		},
	}

	rows := sqlmock.NewRows([]string{"count(1)"}).AddRow(count)
	mock.ExpectQuery("SELECT .+ FROM").
		WillReturnRows(rows)

	iterationCount := count / migration.config.Batch.Size
	mod := count % migration.config.Batch.Size
	if mod > 0 {
		iterationCount++
	}

	for i := int64(0); i < iterationCount; i++ {
		mock.ExpectExec("INSERT INTO").
			WithArgs(migration.config.Batch.Size).
			WillReturnResult(sqlmock.NewResult(migration.config.Batch.Size, migration.config.Batch.Size))
	}

	err = migration.Run()
	if err != nil {
		t.Error(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
