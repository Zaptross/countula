package statistics

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Collect(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate, channelID string, ct database.Turn) {
	var stats []*database.StatisticRow
	db.Where("user_id = ?", ct.UserID).Find(&stats)
	allStatistics := database.GetAllStatistics()

	if len(stats) == 0 {
		stats = make([]*database.StatisticRow, len(allStatistics))
		for i, stat := range allStatistics {
			stats[i] = &database.StatisticRow{
				UserID: ct.UserID,
				Stat:   stat.Key,
				Value:  0,
			}
		}
	}

	userStats := fromStatisticRowArray(stats)

	// explicitly don't store "games" stat, instead we'll just count
	// the number of rows in the turns table on demand

	userStats[database.StatTurns.Key]++ // turns

	if ct.Correct {
		userStats[database.StatCorrect.Key]++ // correct
		userStats[database.StatStreak.Key]++  // streak (current)
		if userStats[database.StatStreak.Key] > userStats[database.StatStreakMax.Key] {
			userStats[database.StatStreakMax.Key] = userStats[database.StatStreak.Key] // streak (max)
		}
	} else {
		userStats[database.StatIncorrect.Key]++ // incorrect
		userStats[database.StatStreak.Key] = 0  // streak (current)
		if ct.Turn > userStats[database.StatGameStreakBreak.Key] {
			userStats[database.StatGameStreakBreak.Key] = ct.Turn // game streak break
		}
	}

	stats = updateStatisticRowArray(stats, userStats) // apply updates to rows

	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "stat"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&stats).Error

	if err != nil {
		s.ChannelMessageSend(channelID, "Error updating statistics.")
		println(fmt.Sprintf("%s - Failed to update statistics: %s", time.Now().Format(time.RFC3339), err.Error()))

		return
	}
}
