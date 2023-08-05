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
	StatTurns           = makeStatistic("turns", "Total number of turns")
	StatCorrect         = makeStatistic("correct", "Total number of correct answers")
	StatIncorrect       = makeStatistic("incorrect", "Total number of incorrect answers")
	StatStreak          = makeStatistic("streak", "Current streak of correct answers")
	StatStreakMax       = makeStatistic("streak_max", "Maximum streak of correct answers")
	StatGames           = makeStatistic("games", "Total number of games participated in")
	StatGameStreakBreak = makeStatistic("game_streak_break", "Longest game you've ended by getting an answer wrong")
)

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
