package handler

import (
	"fmt"
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/esc-chula/esc-discord-bot/internal/instance"
)

type messageHandler struct {
}

func NewMessageHandler() *messageHandler {
	return &messageHandler{}
}

func (h *messageHandler) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var err error

	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Content) == 10 {
		validStudentId, _ := regexp.MatchString(`^6\d{7}21$`, m.Content)
		if validStudentId {
			_, err = s.ChannelMessageSend(m.ChannelID, "🔍 กำลังค้นหาข้อมูล...")
			if err != nil {
				log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
				return
			}

			usersData := instance.GetUsersInstance()
			for _, user := range usersData {
				if user["StudentId"] == m.Content {
					_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("คุณคือ \n\n👤 **%v** \n\n✅  หากข้อมูลถูกต้องให้พิมพ์ `confirm` \n❌  หากข้อมูลผิดให้พิมพ์ `cancel`", user["FullName"]))
					if err != nil {
						log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
						return
					}
					return
				}
			}

			_, err = s.ChannelMessageSend(m.ChannelID, "❌  **ไม่พบข้อมูลของคุณ**\nกรุณาตรวจสอบว่ากรอก Contact List แล้วหรือยัง (https://intania.link/esc67-contact-list-form)\nหรือติดต่อฝ่าย TECH ได้เลย")
			if err != nil {
				log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
				return
			}
		} else {
			_, err = s.ChannelMessageSend(m.ChannelID, "❌  รหัสนิสิตไม่ถูกต้อง กรุณากรอกใหม่อีกครั้ง")
			if err != nil {
				log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
				return
			}
		}
	}
}
