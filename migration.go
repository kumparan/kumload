package kumload

import (
	"fmt"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

// Migration :nodoc:
type Migration interface {
	Run() error
}

type migrationImpl struct {
	db     *gorm.DB
	config *Config
}

// NewMigration create new migration
func NewMigration(db *gorm.DB, config *Config) Migration {
	return &migrationImpl{
		db:     db,
		config: config,
	}
}

func (m *migrationImpl) Run() error {
	var count int64
	limit := m.config.Batch.Size

	countScript := fmt.Sprintf(
		"SELECT count(1) FROM %s.%s WHERE %s NOT IN (SELECT %s FROM %s)",
		m.config.Database.Source.Name, m.config.Database.Source.Table,
		m.config.Database.Source.PrimaryKey, m.config.Database.Target.PrimaryKey,
		m.config.Database.Target.Table,
	)

	err := m.db.Raw(countScript).Count(&count).Error
	if err != nil {
		return err
	}

	log.Info("data: ", count)
	log.Info("batch_size: ", m.config.Batch.Size)

	iterationCount := count / limit
	mod := count % limit
	if mod > 0 {
		iterationCount++
	}

	log.Info("iteration: ", iterationCount)

	script := m.config.Script
	if script == "" {
		script = m.generateQuery(m.config)
	}

	progressBar := pb.New(int(iterationCount))
	progressBar.SetWriter(os.Stdout)
	progressBar.SetTemplate(pb.Full)

	progressBar.Start()
	for i := int64(0); i < iterationCount; i++ {
		err := m.db.Exec(script, limit).Error
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(m.config.Batch.Interval) * time.Second)
		progressBar.Increment()

		if i == iterationCount-2 && mod > 0 {
			limit = mod
		}
	}
	progressBar.Finish()

	return nil
}

func (m *migrationImpl) generateQuery(config *Config) string {
	script := `
INSERT INTO %s(
	%s
) 
SELECT
	%s
FROM 
	%s
WHERE 
	%s NOT IN (
		SELECT %s
		FROM %s
	)
ORDER BY %s ASC LIMIT ?`

	var insertParams, fromParams string
	for k, v := range config.Mappings {
		fromParams += fmt.Sprintf(`%s,`, k)
		insertParams += fmt.Sprintf(`"%s",`, v)
	}

	// remove , in the end of params
	insertParams = insertParams[:len(insertParams)-1]
	fromParams = fromParams[:len(fromParams)-1]

	return fmt.Sprintf(script,
		config.Database.Target.Table,
		insertParams,
		fromParams,
		config.Database.Source.Name+"."+config.Database.Source.Table,
		config.Database.Source.PrimaryKey,
		config.Database.Target.PrimaryKey,
		config.Database.Target.Table,
		config.Database.Source.OrderKey,
	)
}
