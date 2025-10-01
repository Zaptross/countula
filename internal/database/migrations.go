package database

import (
	"github.com/zaptross/countula/internal/database/migrations"
	"gorm.io/gorm"
)

type Migration interface {
	Up(*gorm.DB) error
	Down(*gorm.DB) error
}

func RunMigrationsUp(db *gorm.DB) {
	// check if any tables exist, in case of a fresh database
	var tableCount int64
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableCount)

	if tableCount == 0 {
		// no tables exist, so we can skip migrations
		err := db.AutoMigrate(&Turn{}, &AuditLog{}, &StatisticRow{}, &ServerConfig{}, &RuleSetting{}, &ServiceConfig{})
		if err != nil {
			panic(err)
		}
		return
	}

	migrations := []Migration{
		// order matters
		&migrations.AddChannelToStatistic{},
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		for _, m := range migrations {
			if err := m.Up(tx); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Turn{}, &AuditLog{}, &StatisticRow{}, &ServerConfig{}, &RuleSetting{}, &ServiceConfig{})

	if err != nil {
		panic(err)
	}
}
