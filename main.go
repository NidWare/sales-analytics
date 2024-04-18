package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sales-count/config"
	"sales-count/queryBuilder"
	"strings"
	"time"
)

var (
	bot *tgbotapi.BotAPI
	db  *sql.DB
	cfg *config.Config
)

func main() {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal("Error getting executable path:", err)
	}

	// Get the directory of the executable
	execDir := filepath.Dir(execPath)

	// Construct the absolute path to the config.yml file
	configPath := filepath.Join(execDir, "config.yml")

	// Load the configuration using the absolute path
	cfg, err = config.LoadConfig(configPath)
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	bot, err = tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

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

		if !isAdmin(update.Message.Chat.ID) {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				handleStartCommand(update.Message)
			}
		} else if update.Message.Text == "GET REPORT" {
			handleGetReportRequest(update.Message)
		} else {
			handleDateInput(update.Message)
		}
	}
}

func handleStartCommand(message *tgbotapi.Message) {
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

func handleGetReportRequest(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Please enter the start and end dates in the format: dd-mm-yy dd-mm-yy")
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func handleDateInput(message *tgbotapi.Message) {
	dates := strings.Split(message.Text, " ")
	if len(dates) != 2 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid date format. Please try again.")
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return
	}

	startDate, err := time.Parse("02-01-06", dates[0])
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid start date format. Please try again.")
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return
	}
	startDate = startDate.AddDate(0, 0, 1).Add(3 * time.Hour)

	endDate, err := time.Parse("02-01-06", dates[1])
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Invalid end date format. Please try again.")
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return
	}

	endDate = endDate.AddDate(0, 0, 1).Add(3 * time.Hour)

	managerIDs, err := getManagerIDsFromDatabase(db)
	fmt.Println(managerIDs)
	if err != nil {
		log.Println("Error retrieving manager IDs:", err)
		msg := tgbotapi.NewMessage(message.Chat.ID, "An error occurred while retrieving manager IDs.")
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return
	}

	projectIDs := map[string]string{"Anika, Mia": "1206443281308330", "Amy": "1207064733107968"} // Amy | Anika, Mia

	var reportText string
	reportText += "<pre>"
	reportText += fmt.Sprintf("%-20s | %-10s | %10s\n", "Manager Name", "Model", "Sum")
	reportText += fmt.Sprintf("%-20s | %-10s | %10s\n", "--------------------", "----------", "----------")

	for name, projectID := range projectIDs {

		for _, managerID := range managerIDs {
			managerName, err := getManagerNameFromDatabase(db, managerID)
			if err != nil {
				log.Printf("Error retrieving manager name for manager %s: %v\n", managerID, err)
				continue
			}

			sum, err := calculateSumByManagerID(managerID, startDate, endDate, []string{projectID})
			if err != nil {
				log.Printf("Error calculating sum for manager %s and model %s: %v\n", managerID, name, err)
				continue
			}

			reportText += fmt.Sprintf("%-20s | %-10s | $%9.2f\n", managerName, name, sum)
		}
	}

	reportText += "</pre>"

	if reportText == "<pre></pre>" {
		reportText = "No data found for the specified period."
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, reportText)
	msg.ParseMode = "HTML"
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	}

	if reportText == "" {
		reportText = "No data found for the specified period."
	}
}

func getManagerIDsFromDatabase(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT id FROM managers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var managerIDs []string
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		managerIDs = append(managerIDs, id)
	}

	return managerIDs, nil
}

func getManagerNameFromDatabase(db *sql.DB, managerID string) (string, error) {
	var name string
	err := db.QueryRow("SELECT name FROM managers WHERE id = ?", managerID).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

type Task struct {
	Gid          string        `json:"gid"`
	Assignee     Assignee      `json:"assignee"`
	CompletedAt  string        `json:"completed_at"`
	CustomFields []CustomField `json:"custom_fields"`
}

type Assignee struct {
	Gid          string `json:"gid"`
	ResourceType string `json:"resource_type"`
}

type CustomField struct {
	Gid          string  `json:"gid"`
	DisplayValue *string `json:"display_value"`
}

type Response struct {
	Data []Task `json:"data"`
}

func calculateSumByManagerID(managerID string, startDate, endDate time.Time, projectIDs []string) (float64, error) {
	url := buildAsanaQueryURL(managerID, startDate, endDate, projectIDs)

	response, err := makeAsanaRequest(url)
	if err != nil {
		return 0, err
	}

	sum, err := calculateSumFromResponse(response)
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func buildAsanaQueryURL(managerID string, startDate, endDate time.Time, projectIDs []string) string {
	QB := queryBuilder.NewAsanaTaskSearchBuilder("1206405818803094")
	QB.AddField("assignee").AddField("completed_at").AddField("custom_fields").AddField("custom_fields.display_value")
	QB.SetPretty(true)
	QB.SetResourceSubtype("default_task")
	QB.SetCompletedBefore(endDate)
	QB.SetCompletedAfter(startDate)
	QB.SetCompleted(true)
	QB.SetSortBy("modified_at")
	QB.SetSortAscending(false)
	QB.AddAssigneeID(managerID) // manager

	for _, projectID := range projectIDs {
		QB.AddProjectID(projectID)
	}

	return QB.Build()
}

func makeAsanaRequest(url string) (Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Response{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", cfg.AsanaToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Response{}, err
	}

	return response, nil
}

func calculateSumFromResponse(response Response) (float64, error) {
	var sum float64
	for _, task := range response.Data {
		for _, field := range task.CustomFields {
			if field.DisplayValue != nil {
				var value float64
				_, err := fmt.Sscanf(*field.DisplayValue, "%f", &value)
				if err == nil {
					sum += value
				}
			}
		}
	}
	return sum, nil
}

func isAdmin(userID int64) bool {
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
