package handler

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/esc-chula/esc-discord-bot/config"
	"github.com/esc-chula/esc-discord-bot/internal/instance"
	"github.com/esc-chula/esc-discord-bot/internal/repository"
)

type messageHandler struct {
	userRepo *repository.UserRepository
}

func NewMessageHandler(userRepo *repository.UserRepository) *messageHandler {
	return &messageHandler{
		userRepo: userRepo,
	}
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
		if user["Student ID"] == m.Content {
			userData = user
			break
		}
	}

	studentId := m.Content

	validStudentId, _ := regexp.MatchString(`^6\d{7}21$`, studentId)

	// TYPED STUDENT ID is FOUND
	if userData != nil {
		// CHECKING
		if userData["Bot Status"] == "unconfirmed" && len(m.Content) == 10 {
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("คุณคือ \n\n👤 **%v** \n\n✅  หากข้อมูลถูกต้องให้พิมพ์ `confirm` \n❌  หากข้อมูลผิดให้พิมพ์ `cancel`", userData["Full Name"]))
			if err != nil {
				log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
				return
			}

			instance.SetUserDataByStudentId(studentId, "Bot Status", "confirming")
			instance.SetUserDataByStudentId(studentId, "Discord ID", m.Author.ID)

			return
		}

		// CONFIRMING
		if userData["Bot Status"] == "confirming" && userData["Discord ID"] == m.Author.ID && m.Content == "confirm" {
			_, err = s.ChannelMessageSend(m.ChannelID, "💭  กำลังทำการให้ Role กรุณารอสักครู่...")
			if err != nil {
				log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
				return
			}

			roles := cfg.Discord.DepartmentRoles[userData["Department"].(string)]

			if userData["Year"] != nil {
				yearRole := cfg.Discord.YearRoles[userData["Year"].(string)]
				roles = append(roles, yearRole)
			}

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

			err = h.userRepo.UpdateUserDiscordId(userData["Id"], m.Author.ID)
			if err != nil {
				log.Printf("User: %v, Error updating user data: %v", m.Author.ID, err)

				_, err = s.ChannelMessageSend(m.ChannelID, "❌  ไม่สามารถอัพเดทข้อมูลของคุณได้ กรุณาติดต่อฝ่าย TECH")
				if err != nil {
					log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
				}

				return
			}

			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("✅  คุณได้รับ Role: `%v`\n\nสามารถกลับไปที่ดิสคอร์ด %v ได้เลย!\nหรือหากติดปัญหาอะไรฝ่าย TECH พร้อมให้ความช่วยเหลือคับบ!", strings.Join(rolesName, "`, `"), cfg.Discord.ServerName))
			if err != nil {
				log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
				return
			}

			instance.SetUserDataByDiscordId(m.Author.ID, "Bot Status", "confirmed")

			return
		}

		// CANCELING
		if userData["Bot Status"] == "confirming" && userData["Discord ID"] == m.Author.ID && m.Content == "cancel" {
			instance.SetUserDataByDiscordId(m.Author.ID, "Bot Status", "unconfirmed")
			instance.SetUserDataByDiscordId(m.Author.ID, "Discord ID", nil)

			_, err = s.ChannelMessageSend(m.ChannelID, "❌  ยกเลิกการยืนยันข้อมูลแล้ว\nสามารถพิมพ์รหัสนิสิตอีกครั้งเพื่อค้นหาและยืนยันการให้ Role")
			if err != nil {
				log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
				return
			}

			return
		}

		// CORRECTING USER CONFIRMATION
		if userData["Bot Status"] == "confirming" {
			_, err = s.ChannelMessageSend(m.ChannelID, "กรุณาพิมพ์ `confirm` หรือ `cancel` เพิ่อดำเนินการต่อ")
			if err != nil {
				log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
				return
			}
			return
		}

		// ALREADY CONFIRMED
		if userData["Bot Status"] == "confirmed" {
			_, err = s.ChannelMessageSend(m.ChannelID, "👋  รหัสนิสิตนี้ได้รับ Role แล้ว สามารถกลับไปที่ดิสคอร์ดได้เลย!\nหรือหากติดปัญหาอะไร สามารถติดต่อฝ่าย TECH ได้เลยทันที")
			if err != nil {
				log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
				return
			}
		}
	} else {
		// NOT FOUND
		if validStudentId {
			_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("❌  **ไม่พบข้อมูลของคุณ**\nกรุณาตรวจสอบว่ากรอก Contact List แล้วหรือยัง (%v)\nหรือติดต่อฝ่าย TECH ได้เลย", cfg.Discord.ContactListForm))
			if err != nil {
				log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
				return
			}

			return
		}

		// UNKNOW COMMAND
		_, err = s.ChannelMessageSend(m.ChannelID, "กรุณาพิมพ์รหัสนิสิตของตัวเองเพื่อรับ Role")
		if err != nil {
			log.Printf("User: %v, Error sending message: %v", m.Author.ID, err)
			return
		}

		return
	}
}
