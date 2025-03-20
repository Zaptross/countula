package migrations

import (
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

type AddChannelToStatistic struct{}

func (m *AddChannelToStatistic) Up(db *gorm.DB) error {
	// Check if the column already exists
	var columnExists bool
	query := `
    SELECT COUNT(*) > 0
    FROM information_schema.columns
    WHERE table_name = 'statistic_rows'
    AND column_name = 'channel_id'
  `

	err := db.Raw(query).Scan(&columnExists).Error
	if err != nil {
		return fmt.Errorf("failed to check if column exists: %w", err)
	}

	if !columnExists {
		// Add channel_id column
		err = db.Exec(`ALTER TABLE statistic_rows ADD COLUMN channel_id TEXT`).Error
		if err != nil {
			return fmt.Errorf("failed to add channel_id column: %w", err)
		}

		// Set default value for existing records
		err = db.Exec(`UPDATE statistic_rows SET channel_id = 'global' WHERE channel_id IS NULL`).Error
		if err != nil {
			return fmt.Errorf("failed to set default channel_id: %w", err)
		}

		// Drop the current primary key, if it exists
		err = db.Exec(`ALTER TABLE statistic_rows DROP CONSTRAINT statistic_rows_pkey`).Error

		if err != nil {
			return fmt.Errorf("failed to drop primary key: %w", err)
		}

		// Add a new composite primary key
		err = db.Exec(`ALTER TABLE statistic_rows ADD PRIMARY KEY (user_id, channel_id, stat)`).Error

		if err != nil {
			return fmt.Errorf("failed to add composite primary key: %w", err)
		}

		slog.Info("Migration AddChannelToStatistic completed successfully")
	} else {
		slog.Debug("Skipping migration: AddChannelToStatistic")
	}

	return nil
}

func (m *AddChannelToStatistic) Down(db *gorm.DB) error {
	// Drop the column if it exists
	var columnExists bool
	query := `
    SELECT COUNT(*) > 0
    FROM information_schema.columns
    WHERE table_name = 'statistic_rows'
    AND column_name = 'channel_id'
  `

	err := db.Raw(query).Scan(&columnExists).Error
	if err != nil {
		return fmt.Errorf("failed to check if column exists: %w", err)
	}

	if columnExists {
		err = db.Exec(`ALTER TABLE statistic_rows DROP COLUMN channel_id`).Error
		if err != nil {
			return fmt.Errorf("failed to drop channel_id column: %w", err)
		}
	}

	slog.Info("Migration AddChannelToStatistic reverted successfully")

	return nil
}
