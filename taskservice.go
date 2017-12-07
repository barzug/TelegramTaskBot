package main

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
)

type User struct {
	ChatId   int64
	UserName string
}

type Task struct {
	Creator  User
	Executor User
	Text     string
	Id       int64
}

type TasksService struct {
	tasks []Task
	id    int64
	mu    sync.Mutex
}

func (task *Task) taskParse(user User, withAssignee bool) string {
	resultString := strconv.FormatInt(task.Id, 10) + `. ` + task.Text + ` by @` + task.Creator.UserName

	switch task.Executor.UserName {
	case ``:
		resultString += "\n/assign_" + strconv.FormatInt(task.Id, 10)
	case user.UserName:
		if withAssignee {
			resultString += "\nassignee: я"
		}
		resultString += "\n/unassign_" + strconv.FormatInt(task.Id, 10) + " /resolve_" + strconv.FormatInt(task.Id, 10)
	default:
		resultString += "\nassignee: @" + task.Executor.UserName
	}
	return resultString
}

func (t *TasksService) addTask(creator User, text string) int64 {
	newTask := Task{}

	newTask.Id = atomic.AddInt64(&t.id, 1)

	newTask.Creator = creator
	newTask.Text = text

	t.tasks = append(t.tasks, newTask)

	return newTask.Id
}

func (t *TasksService) getTasks() []Task {
	return t.tasks
}

func (t *TasksService) getTasksByCreator(creator User) []Task {
	var tasksByCreator []Task

	for _, task := range t.tasks {
		if task.Creator == creator {
			tasksByCreator = append(tasksByCreator, task)
		}
	}

	return tasksByCreator
}

func (t *TasksService) getTasksByExecutor(executor User) []Task {
	var tasksByExecutor []Task

	for _, task := range t.tasks {
		if task.Executor == executor {
			tasksByExecutor = append(tasksByExecutor, task)
		}
	}

	return tasksByExecutor

}

func (t *TasksService) getTaskById(id int64) (*Task, error) {
	for i := range t.tasks {
		if t.tasks[i].Id == id {
			return &t.tasks[i], nil
		}
	}
	return nil, errors.New("Нет задачи с таким id")
}

func (t *TasksService) removeTaskById(id int64) error {
	for i := range t.tasks {
		if t.tasks[i].Id == id {
			t.mu.Lock()
			t.tasks = append(t.tasks[:i], t.tasks[i+1:]...)
			t.mu.Unlock()
			return nil
		}
	}
	return errors.New("Нет задачи с таким id")
}

func (t *TasksService) assignTask(executor User, id int64) (string, User, error) {
	task, err := t.getTaskById(id)
	if err != nil {
		return "", User{}, err
	}

	var userToNotify User

	if task.Executor.UserName == "" {
		userToNotify = task.Creator
	} else {
		userToNotify = task.Executor
	}

	task.Executor = executor

	return task.Text, userToNotify, nil
}

func (t *TasksService) unassignTask(executor User, id int64) (string, User, error) {
	task, err := t.getTaskById(id)
	if err != nil {
		return "", User{}, err
	}

	if task.Executor.UserName != executor.UserName {
		return "", User{}, errors.New("Задача не на вас")
	}

	userToNotify := task.Creator

	task.Executor = User{}

	return task.Text, userToNotify, nil
}

func (t *TasksService) resolveTask(id int64) (string, User, error) {
	task, err := t.getTaskById(id)
	if err != nil {
		return "", User{}, err
	}

	userToNotify := task.Creator

	err = t.removeTaskById(id)
	if err != nil {
		return "", User{}, err
	}

	return task.Text, userToNotify, nil
}
