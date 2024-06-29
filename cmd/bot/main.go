package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/esc-chula/esc-discord-bot/config"
	"github.com/esc-chula/esc-discord-bot/internal/handler"
	"github.com/esc-chula/esc-discord-bot/internal/instance"
	"github.com/esc-chula/esc-discord-bot/pkg/utils"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting bot...")

	// LOAD CONFIG
	configPath := config.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	config.SetInstance(cfg)

	usersData, err := utils.GetUsersData()
	if err != nil {
		log.Fatalf("GetUsersData: %v", err)
	}
	instance.SetUsersInstance(usersData)

	// HANDLER
	guildMemberHandler := handler.NewGuildMemberHandler()
	messageHandler := handler.NewMessageHandler()
	webhookhandler := handler.NewWebhookHandler()

	// BOT
	dg, err := discordgo.New("Bot " + cfg.Bot.Token)
	if err != nil {
		log.Fatalf("Error creating Discord session, %v", err)
		return
	}

	dg.AddHandler(guildMemberHandler.Add)
	dg.AddHandler(messageHandler.MessageCreate)

	dg.Identify.Intents = discordgo.IntentsAll

	// WEBHOOK SERVER
	router := mux.NewRouter()
	router.HandleFunc("/webhook", webhookhandler.Received).Methods("POST")

	go func() {
		http.Handle("/", router)
		log.Printf("Webhook server is listening on port %s", cfg.Webhook.Port)
		log.Fatal(http.ListenAndServe(":"+cfg.Webhook.Port, nil))
	}()

	utils.SigHandler(dg)
}
