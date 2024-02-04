package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	TELEGRAM_API_URL = "https://api.telegram.org/bot%s/sendMessage"

	// Driunk reminder
	DRINK_REMINDER = "*Ayo Minum!* 🥛\nMinimal 1 gelas"
)

type telegramMessage struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

var err error
var message string
var message_type *string
var chatIDs []string

func sendTelegramNotification(botToken string, chatID int64, message string) error {
	// Create the Telegram message payload
	telegramMsg := telegramMessage{
		ChatID:    chatID,
		Text:      message,
		ParseMode: "Markdown",
	}

	// Convert the message payload to a JSON string
	msgBytes, err := json.Marshal(telegramMsg)
	if err != nil {
		return err
	}

	// Create a new HTTP request to send the Telegram message
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(TELEGRAM_API_URL, botToken), bytes.NewReader(msgBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request to the Telegram API
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the response status code indicates success
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get bot token and chat ID from environment variables
	botToken := os.Getenv("BOT_TOKEN")
	chatIDsEnv := os.Getenv("CHAT_ID")

	// Split the string on the comma delimiter
	chatIDs = strings.Split(chatIDsEnv, ",")

	message_type = flag.String("type", "", "message type is drink_reminder")

	flag.Parse()

	if len(*message_type) == 0 {
		log.Fatal("Type is empty!")
	}

	switch *message_type {
	case "drink_reminder":
		message = DRINK_REMINDER
	default:
		log.Fatal("Type not found!")
	}

	for _, chatID := range chatIDs {
		// Parse chat ID from string to int64
		chatIDint, err := strconv.ParseInt(chatID, 10, 64)
		if err != nil {
			log.Fatal("Error parsing chat ID")
		}

		// Send a notification message to the specified chat ID
		err = sendTelegramNotification(botToken, chatIDint, message)
		if err != nil {
			fmt.Println(err)
		}
	}
}
