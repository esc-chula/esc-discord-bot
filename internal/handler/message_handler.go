package handler

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/esc-chula/esc-discord-bot/config"
	"github.com/esc-chula/esc-discord-bot/internal/instance"
)

type messageHandler struct {
}

func NewMessageHandler() *messageHandler {
	return &messageHandler{}
}

func (h *messageHandler) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	var err error

	cfg := config.GetConfig()

	if m.Author.ID == s.State.User.ID || m.GuildID != "" {
		return
	}

	log.Printf("Message from %v: %v", m.Author.ID, m.Content)

	usersData := instance.GetUsersInstance()

	var userData map[string]interface{}

	for _, user := range usersData {
		if user["Discord ID"] == m.Author.ID {
			userData = user
			break
		}
		if user["Student Id"] == m.Content {
			userData = user
			break
		}
	}

	if userData == nil {
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌  **ไม่พบข้อมูลของคุณ**\nกรุณาตรวจสอบว่ากรอก Contact List แล้วหรือยัง (%v)\nหรือติดต่อฝ่าย TECH ได้เลย", cfg.Discord.ContactListForm))
		if err != nil {
			log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
			return
		}
		return
	}

	if len(m.Content) == 10 && userData["Bot Status"] == "unconfirmed" {
		studentId := m.Content

		validStudentId, _ := regexp.MatchString(`^6\d{7}21$`, studentId)
		if !validStudentId {
			_, err = s.ChannelMessageSend(m.ChannelID, "❌  รหัสนิสิตไม่ถูกต้อง กรุณากรอกใหม่อีกครั้ง")
			if err != nil {
				log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
				return
			}
		}

		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("คุณคือ \n\n👤 **%v** \n\n✅  หากข้อมูลถูกต้องให้พิมพ์ `confirm` \n❌  หากข้อมูลผิดให้พิมพ์ `cancel`", userData["Full Name"]))
		if err != nil {
			log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
			return
		}

		instance.SetUserDataByStudentId(studentId, "Bot Status", "confirming")
		instance.SetUserDataByStudentId(studentId, "Discord ID", m.Author.ID)

		return
	}

	if m.Content == "confirm" && userData["Discord ID"] == m.Author.ID && userData["Bot Status"] == "confirming" {
		instance.SetUserDataByDiscordId(m.Author.ID, "Bot Status", "confirmed")

		_, err = s.ChannelMessageSend(m.ChannelID, "💭  กำลังทำการให้ Role กรุณารอสักครู่...")
		if err != nil {
			log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
			return
		}

		departmentRoles := cfg.Discord.DepartmentRoles[userData["Department"].(string)]
		yearRole := cfg.Discord.YearRoles[userData["Year"].(string)]

		roles := append(departmentRoles, yearRole)

		for _, role := range roles {
			err = s.GuildMemberRoleAdd(cfg.Discord.ServerId, m.Author.ID, role)
			if err != nil {
				log.Printf("User: %v, Error adding role: %v", m.Author.ID, err)
				return
			}
		}

		rolesName := []string{}
		for _, role := range roles {
			serverRole, err := s.State.Role(cfg.Discord.ServerId, role)
			if err != nil {
				continue
			}
			rolesName = append(rolesName, serverRole.Name)
		}

		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅  คุณได้รับ Role: `%v`\n\nสามารถกลับไปที่ดิสคอร์ด %v ได้เลย!\nหรือหากติดปัญหาอะไรฝ่าย TECH พร้อมให้ความช่วยเหลือคับ", strings.Join(rolesName, "`, `"), cfg.Discord.ServerName))
		if err != nil {
			log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
			return
		}

		// TODO: patch usersData to NocoDB

		return
	}

	if m.Content == "cancel" && userData["Discord ID"] == m.Author.ID && userData["Bot Status"] == "confirming" {
		instance.SetUserDataByDiscordId(m.Author.ID, "Bot Status", "unconfirmed")
		instance.SetUserDataByDiscordId(m.Author.ID, "Discord ID", nil)

		_, err = s.ChannelMessageSend(m.ChannelID, "❌  ยกเลิกการยืนยันข้อมูลแล้ว\nสามารถพิมพ์รหัสนิสิตอีกครั้งเพื่อค้นหาและยืนยันการให้ Role")
		if err != nil {
			log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
			return
		}

		return
	}

	if userData["Bot Status"] == "confirming" {
		_, err = s.ChannelMessageSend(m.ChannelID, "กรุณาพิมพ์ `confirm` หรือ `cancel` เพิ่อดำเนินการต่อ")
		if err != nil {
			log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
			return
		}
		return
	}

	if userData["Bot Status"] == "confirmed" {
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, "กรุณาพิมพ์รหัสนิสิตของตัวเองเพื่อรับ Role")
	if err != nil {
		log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
		return
	}
}
