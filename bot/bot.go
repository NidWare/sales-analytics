package bot

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"sales-count/asana"
	"sales-count/config"
	"sales-count/database"
	"strings"
	"time"
)

func Start(cfg *config.Config, db *sql.DB) {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !isAdmin(update.Message.Chat.ID, cfg) {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				handleStartCommand(bot, update.Message)
			}
		} else if update.Message.Text == "GET REPORT" {
			handleGetReportRequest(bot, update.Message)
		} else {
			handleDateInput(bot, update.Message, db, cfg)
		}
	}
}

func handleStartCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome to the Asana Report Bot!")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("GET REPORT"),
		),
	)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func handleGetReportRequest(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Please enter the start and end dates in the format:\ndd-mm-yy dd-mm-yy")
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func handleDateInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message, db *sql.DB, cfg *config.Config) {
	dates := strings.Split(message.Text, " ")
	if len(dates) != 2 {
		sendErrorMessage(bot, message.Chat.ID, "Invalid date format. Please try again.")
		return
	}

	startDate, err := time.Parse("02-01-06", dates[0])
	if err != nil {
		sendErrorMessage(bot, message.Chat.ID, "Invalid start date format. Please try again.")
		return
	}
	startDate = startDate.Add(3 * time.Hour)

	endDate, err := time.Parse("02-01-06", dates[1])
	if err != nil {
		sendErrorMessage(bot, message.Chat.ID, "Invalid end date format. Please try again.")
		return
	}
	endDate = endDate.Add(3 * time.Hour)

	managerIDs, err := database.GetManagerIDsFromDatabase(db)
	if err != nil {
		log.Println("Error retrieving manager IDs:", err)
		sendErrorMessage(bot, message.Chat.ID, "An error occurred while retrieving manager IDs.")
		return
	}

	projectIDs := map[string]string{"Anika, Mia": "1206443281308330", "Amy": "1207064733107968"}

	reportText := generateReportText(db, managerIDs, projectIDs, startDate, endDate, cfg)

	msg := tgbotapi.NewMessage(message.Chat.ID, reportText)
	msg.ParseMode = "markdown"
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func generateReportText(db *sql.DB, managerIDs []string, projectIDs map[string]string, startDate, endDate time.Time, cfg *config.Config) string {
	var reportText string
	for name, projectID := range projectIDs {
		reportText += fmt.Sprintf("*%s:*\n", name)
		for _, managerID := range managerIDs {
			managerName, err := database.GetManagerNameFromDatabase(db, managerID)
			if err != nil {
				log.Printf("Error retrieving manager name for manager %s: %v\n", managerID, err)
				continue
			}

			sum, err := asana.CalculateSumByManagerID(managerID, startDate, endDate, []string{projectID}, cfg.AsanaToken)
			if err != nil {
				log.Printf("Error calculating sum for manager %s and model %s: %v\n", managerID, name, err)
				continue
			}

			reportText += fmt.Sprintf("%s: $%.2f\n", managerName, sum)
		}
		reportText += "\n\n"
	}

	if reportText == "" {
		reportText = "No data found for the specified period."
	}

	return reportText
}

func sendErrorMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func isAdmin(userID int64, cfg *config.Config) bool {
	if cfg == nil {
		log.Println("Configuration is not initialized")
		return false
	}

	for _, adminID := range cfg.AdminIDs {
		if userID == adminID {
			return true
		}
	}
	return false
}
