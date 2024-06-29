package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Bot     BotConfig
	Webhook WebhookConfig
	Discord DiscordConfig
}

type BotConfig struct {
	Version string
	Token   string // env: BOT_TOKEN
}

type WebhookConfig struct {
	Port              string
	Secret            string // env: WEBHOOK_SECRET
	NocoDBAPIEndpoint string
	NocoDBAPIToken    string // env: WEBHOOK_NOCODB_API_TOKEN
	NocoDBTableId     string
}

type DiscordConfig struct {
	ServerId         string
	WelcomeChannelId string
	Roles            map[string]string
}

func overrideWithEnv(c *Config) {
	if value, exists := os.LookupEnv("BOT_TOKEN"); exists {
		c.Bot.Token = value
	}

	if value, exists := os.LookupEnv("WEBHOOK_SECRET"); exists {
		c.Webhook.Secret = value
	}
	if value, exists := os.LookupEnv("WEBHOOK_NOCODB_API_TOKEN"); exists {
		c.Webhook.NocoDBAPIToken = value
	}
}

func LoadConfig(filename string) (*viper.Viper, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	overrideWithEnv(&c)

	return &c, nil
}

func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./config/config-local"
}
