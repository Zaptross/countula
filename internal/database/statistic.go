package database

import (
	"time"

	"gorm.io/gorm"
)

type StatisticRow struct {
	UserID    string `gorm:"primaryKey"`
	Stat      string `gorm:"primaryKey"`
	Value     int
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Statistic struct {
	Key         string
	Description string
}

var (
	StatTurns           = makeStatistic("turns", "Turns played")
	StatCorrect         = makeStatistic("correct", "Your correct guesses")
	StatIncorrect       = makeStatistic("incorrect", "Your fuckups")
	StatStreak          = makeStatistic("streak", "Current streak")
	StatStreakMax       = makeStatistic("streak_max", "Your longest streak")
	StatGameStreakBreak = makeStatistic("game_streak_break", "Longest game you've ended by getting an answer wrong")
)

var allStatistics = []*Statistic{
	StatTurns,
	StatCorrect,
	StatIncorrect,
	StatStreak,
	StatStreakMax,
	StatGameStreakBreak,
}

func GetAllStatistics() []*Statistic {
	return allStatistics
}

func GetStatByKey(key string) *Statistic {
	for _, stat := range allStatistics {
		if stat.Key == key {
			return stat
		}
	}
	return nil
}

func makeStatistic(key string, description string) *Statistic {
	return &Statistic{
		Key:         key,
		Description: description,
	}
}

func (s *Statistic) ToRow(userID string, value int) StatisticRow {
	return StatisticRow{
		UserID: userID,
		Stat:   s.Key,
		Value:  value,
	}
}

func (s *Statistic) ForUser(db *gorm.DB, userID string) StatisticRow {
	var row StatisticRow
	db.Where("user_id = ? AND stat = ?", userID, s.Key).First(&row)
	return row
}
