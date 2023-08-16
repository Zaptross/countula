package statistics

import "github.com/zaptross/countula/internal/database"

type UserStats = map[string]int

func fromStatisticRowArray(rows []*database.StatisticRow) UserStats {
	stats := make(UserStats)
	for _, row := range rows {
		stats[row.Stat] = row.Value
	}
	return stats
}

func updateStatisticRowArray(rows []*database.StatisticRow, stats UserStats) []*database.StatisticRow {
	for _, row := range rows {
		if val, ok := stats[row.Stat]; ok {
			row.Value = val
		}
	}
	return rows
}
