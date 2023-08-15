package statistics

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/countula/internal/database"
	"gorm.io/gorm"
)

func Display(db *gorm.DB, s *discordgo.Session, m *discordgo.MessageCreate) {
	var stats []database.StatisticRow
	db.Where("user_id = ?", m.Author.ID).Find(&stats)

	reply := ""

	for _, stat := range stats {
		reply += fmt.Sprintf("%s: %d\n", database.GetStatByKey(stat.Stat).Description, stat.Value)
	}

	s.ChannelMessageSendReply(m.ChannelID, reply, m.Message.Reference())
}
