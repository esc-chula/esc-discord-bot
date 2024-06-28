package utils

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func SigHandler(dg *discordgo.Session) {
	err := dg.Open()
	if err != nil {
		log.Println("Error opening connection,", err)
		return
	}

	log.Println("Bot is now running...")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-sigCh

	log.Println("Shutting down bot...")
	dg.Close()
}
