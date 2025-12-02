package database

import (
	"time"

	"gorm.io/gorm"
)

type Migration interface {
	ID() int
	Up(*gorm.DB) error
	Down(*gorm.DB) error
}

type MigrationRun struct {
	MigrationID int `gorm:"primaryKey"`
	Run         time.Time
}

func RunMigrationsUp(db *gorm.DB) {
	var err error
	// check if any tables exist, in case of a fresh database
	var tableCount int64
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableCount)

	if tableCount == 0 {
		// no tables exist, so we can skip migrations
		err = migrateTables(db)
		if err != nil {
			panic(err)
		}
		return
	}

	// ensure the migrations table exists
	err = db.AutoMigrate(&MigrationRun{})

	if err != nil {
		panic(err)
	}

	migrations := []Migration{
		// order matters
		&addChannelToStatistic{},
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		for _, m := range migrations {
			// check if this migration has already been run
			var count int64
			tx.Model(&MigrationRun{}).Where("migration_id = ?", m.ID()).Count(&count)

			if count > 0 {
				continue
			}

			if err := m.Up(tx); err != nil {
				return err
			}

			// record that the migration has been run
			migrationRun := MigrationRun{
				MigrationID: m.ID(),
				Run:         time.Now(),
			}

			if err := tx.Create(&migrationRun).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	err = migrateTables(db)

	if err != nil {
		panic(err)
	}
}

func migrateTables(db *gorm.DB) error {
	return db.AutoMigrate(&Turn{}, &AuditLog{}, &StatisticRow{}, &ServerConfig{}, &RuleSetting{}, &ServiceConfig{})
}
