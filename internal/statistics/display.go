package statistics

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

func Display(db *gorm.DB, authorID, channelID string) string {
	var stats []database.StatisticRow
	db.Where("user_id = ? and channel_id in (?, 'global')", authorID, channelID).Find(&stats)

	reply := ""

	for _, stat := range database.GetAllStatistics() {
		statValue := lo.Reduce(stats, func(acc int, s database.StatisticRow, _ int) int {
			if s.Stat == stat.Key {
				return s.Value + acc
			}
			return acc
		}, 0)

		reply += fmt.Sprintf("%s: %d\n", stat.Description, statValue)
	}

	return reply
}

func DisplayGlobal(db *gorm.DB, authorID string) string {
	var stats []database.StatisticRow
	db.Where("user_id = ?", authorID).Find(&stats)

	reply := ""

	for _, stat := range database.GetAllStatistics() {
		statValue := lo.Reduce(stats, func(acc int, s database.StatisticRow, _ int) int {
			if s.Stat == stat.Key {
				return s.Value + acc
			}
			return acc
		}, 0)

		reply += fmt.Sprintf("%s: %d\n", stat.Description, statValue)
	}

	return reply
}
