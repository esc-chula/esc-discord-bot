package handler

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/esc-chula/esc-discord-bot/config"
)

type guildMemberHandler struct {
}

func NewGuildMemberHandler() *guildMemberHandler {
	return &guildMemberHandler{}
}

func (h *guildMemberHandler) Add(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	var err error

	cfg := config.GetConfig()

	if m.GuildID != cfg.Discord.ServerId {
		log.Printf("Server ID mismatch: %v", m.GuildID)
		return
	}

	user := m.User
	channel, err := s.UserChannelCreate(user.ID)
	if err != nil {
		log.Printf("User: %v, Error creating channel: %v", user.ID, err)
		return
	}

	_, err = s.ChannelMessageSend(channel.ID, fmt.Sprintf("%s สามารถพิมพ์รหัสนิสิตของตัวเองมาในแชทเพื่อรับ Role ได้เลย!", user.Mention()))
	if err != nil {
		log.Printf("User: %v, Error sending message: %v", user.ID, err)
	}

	_, err = s.ChannelMessageSend(cfg.Discord.WelcomeChannelId, fmt.Sprintf("**ยินดีต้อนรับ %v สู่ดิสคอร์ด %v 🎉**\nสามารถขอรับ Role ผ่านทาง DM ที่บอทส่งให้ได้เลย \nอย่าลืมเปลี่ยนชื่อตัวเองเป็นชื่อเล่นด้วยหละ!", user.Mention(), cfg.Discord.ServerName))
	if err != nil {
		log.Printf("User: %v, Error sending message: %v", user.ID, err)
	}
}
