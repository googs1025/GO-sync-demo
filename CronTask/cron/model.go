package cron

import (
	"golanglearning/new_project/for-gong-zhong-hao/Concurrent-practice/CronTask"
	"log"
	"os"
)

type TaskScheduler struct {
	tasks []CronTask.TaskInterface
	swap []CronTask.TaskInterface
	add chan CronTask.TaskInterface
	remove chan string
	stop chan struct{}
	Logger CronTask.TaskLogInterface
	lock bool

}

type OnceCron struct {
	*TaskScheduler
}

func NewScheduler() *TaskScheduler {
	return &TaskScheduler{
		tasks: make([]CronTask.TaskInterface, 0),
		swap: make([]CronTask.TaskInterface, 0),
		add: make(chan CronTask.TaskInterface),
		stop: make(chan struct{}),
		remove: make(chan string),
		Logger: log.New(os.Stdout, "[Control]:", log.Ldate|log.Ltime|log.Lshortfile),
	}

}

func NewCron() *OnceCron {
	return &OnceCron{
		TaskScheduler: NewScheduler(),
	}
}
