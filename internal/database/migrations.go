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
