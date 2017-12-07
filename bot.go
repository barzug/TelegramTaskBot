package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

var (
	// @BotFather gives you this
	BotToken   = "496482011:AAFYJnPB6lHpFGv0dl61vtU-CkGnRAEzapw"
	WebhookURL = "https://lolkekbot.herokuapp.com"
)

const (
	getTaskPrefix           = "/tasks"
	addTaskPrefix           = "/new "
	assignTaskPrefix        = "/assign_"
	unassignTaskPrefix      = "/unassign_"
	resolveTaskPrefix       = "/resolve_"
	getTaskByCreatorPrefix  = "/owner"
	getTaskByExecutorPrefix = "/my"
)

func startTaskBot(ctx context.Context, port string) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}

	// bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		panic(err)
	}

	updates := bot.ListenForWebhook("/")

	go http.ListenAndServe(":"+port, nil)
	fmt.Println("start listen :" + port)

	t := new(TasksService)

	for update := range updates {

		user := User{}
		if update.Message == nil || update.Message.Chat == nil {
			continue
		}

		user.ChatId = update.Message.Chat.ID
		user.UserName = update.Message.From.UserName

		switch {
		case strings.HasPrefix(update.Message.Text, getTaskPrefix):
			taskHandler(bot, t, user)

		case strings.HasPrefix(update.Message.Text, getTaskByCreatorPrefix):
			taskByCreatorHandler(bot, t, user)

		case strings.HasPrefix(update.Message.Text, getTaskByExecutorPrefix):
			taskByExecutorHandler(bot, t, user)

		case strings.HasPrefix(update.Message.Text, addTaskPrefix):
			addTaskHandler(bot, t, update.Message.Text, user)

		case strings.HasPrefix(update.Message.Text, assignTaskPrefix):
			assignTaskHandler(bot, t, update.Message.Text, user)

		case strings.HasPrefix(update.Message.Text, unassignTaskPrefix):
			unassignTaskHandler(bot, t, update.Message.Text, user)

		case strings.HasPrefix(update.Message.Text, resolveTaskPrefix):
			resolveTaskHandler(bot, t, update.Message.Text, user)

		default:
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Неизвестная команда",
			))
		}
	}
	return nil
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	err := startTaskBot(context.Background(), port)
	if err != nil {
		panic(err)
	}
}
