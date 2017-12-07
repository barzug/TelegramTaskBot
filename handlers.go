package main

import (
	"strconv"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

func taskHandler(bot *tgbotapi.BotAPI, t *TasksService, requestUser User) {
	var responseMessage string

	tasks := t.getTasks()

	var parsedTasks []string

	if len(tasks) == 0 {
		responseMessage = "Нет задач"
	} else {
		for _, task := range tasks {
			parsedTasks = append(parsedTasks, task.taskParse(requestUser, true))
		}
		responseMessage = strings.Join(parsedTasks, "\n\n")
	}

	bot.Send(tgbotapi.NewMessage(
		requestUser.ChatId,
		responseMessage,
	))
}

func taskByCreatorHandler(bot *tgbotapi.BotAPI, t *TasksService, requestUser User) {
	var responseMessage string

	tasks := t.getTasksByCreator(requestUser)

	var parsedTasks []string

	if len(tasks) == 0 {
		responseMessage = "Нет задач"
	} else {
		for _, task := range tasks {
			parsedTasks = append(parsedTasks, task.taskParse(requestUser, false))
		}
		responseMessage = strings.Join(parsedTasks, "\n\n")
	}

	bot.Send(tgbotapi.NewMessage(
		requestUser.ChatId,
		responseMessage,
	))
}

func taskByExecutorHandler(bot *tgbotapi.BotAPI, t *TasksService, requestUser User) {
	var responseMessage string

	tasks := t.getTasksByExecutor(requestUser)

	var parsedTasks []string

	if len(tasks) == 0 {
		responseMessage = "Нет задач"
	} else {
		for _, task := range tasks {
			parsedTasks = append(parsedTasks, task.taskParse(requestUser, false))
		}
		responseMessage = strings.Join(parsedTasks, "\n\n")
	}

	bot.Send(tgbotapi.NewMessage(
		requestUser.ChatId,
		responseMessage,
	))
}

func addTaskHandler(bot *tgbotapi.BotAPI, t *TasksService, text string, requestUser User) {
	var responseMessage string

	message := strings.TrimPrefix(text, addTaskPrefix)
	id := strconv.FormatInt(t.addTask(requestUser, message), 10)

	responseMessage = `Задача "` + message + `" создана, id=` + id

	bot.Send(tgbotapi.NewMessage(
		requestUser.ChatId,
		responseMessage,
	))
}

func assignTaskHandler(bot *tgbotapi.BotAPI, t *TasksService, text string, requestUser User) {
	var responseMessage string

	message := strings.TrimPrefix(text, assignTaskPrefix)

	id, err := strconv.ParseInt(message, 10, 64)
	if err != nil {
		responseMessage = "Неправильный формат запроса"
		bot.Send(tgbotapi.NewMessage(
			requestUser.ChatId,
			responseMessage,
		))
		return
	}

	userToNotify := User{}
	var taskName string
	taskName, userToNotify, err = t.assignTask(requestUser, id)
	if err != nil {
		responseMessage = err.Error()
		bot.Send(tgbotapi.NewMessage(
			requestUser.ChatId,
			responseMessage,
		))
		return
	}

	responseMessage = `Задача "` + taskName + `" назначена на вас`

	bot.Send(tgbotapi.NewMessage(
		requestUser.ChatId,
		responseMessage,
	))

	if userToNotify.UserName != requestUser.UserName {
		bot.Send(tgbotapi.NewMessage(
			userToNotify.ChatId,
			`Задача "`+taskName+`" назначена на @`+requestUser.UserName,
		))
	}
}

func unassignTaskHandler(bot *tgbotapi.BotAPI, t *TasksService, text string, requestUser User) {
	var responseMessage string

	message := strings.TrimPrefix(text, unassignTaskPrefix)

	id, err := strconv.ParseInt(message, 10, 64)
	if err != nil {
		responseMessage = "Неправильный формат запроса"
		bot.Send(tgbotapi.NewMessage(
			requestUser.ChatId,
			responseMessage,
		))
		return
	}

	userToNotify := User{}
	var taskName string
	taskName, userToNotify, err = t.unassignTask(requestUser, id)
	if err != nil {
		responseMessage = err.Error()
		bot.Send(tgbotapi.NewMessage(
			requestUser.ChatId,
			responseMessage,
		))
		return
	}

	bot.Send(tgbotapi.NewMessage(
		requestUser.ChatId,
		`Принято`,
	))

	if userToNotify.UserName != requestUser.UserName {
		bot.Send(tgbotapi.NewMessage(
			userToNotify.ChatId,
			`Задача "`+taskName+`" осталась без исполнителя`,
		))
	}
}

func resolveTaskHandler(bot *tgbotapi.BotAPI, t *TasksService, text string, requestUser User) {
	var responseMessage string

	message := strings.TrimPrefix(text, resolveTaskPrefix)

	id, err := strconv.ParseInt(message, 10, 64)
	if err != nil {
		responseMessage = "Неправильный формат запроса"
		bot.Send(tgbotapi.NewMessage(
			requestUser.ChatId,
			responseMessage,
		))
		return
	}

	userToNotify := User{}
	var taskName string
	taskName, userToNotify, err = t.resolveTask(id)
	if err != nil {
		responseMessage = err.Error()
		bot.Send(tgbotapi.NewMessage(
			requestUser.ChatId,
			responseMessage,
		))
		return
	}

	bot.Send(tgbotapi.NewMessage(
		requestUser.ChatId,
		`Задача "`+taskName+`" выполнена`,
	))

	if userToNotify.UserName != requestUser.UserName {
		bot.Send(tgbotapi.NewMessage(
			userToNotify.ChatId,
			`Задача "`+taskName+`" выполнена @`+requestUser.UserName,
		))
	}
}
