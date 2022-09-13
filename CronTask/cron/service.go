package cron

import (
	"github.com/google/uuid"
	"golanglearning/new_project/for-gong-zhong-hao/Concurrent-practice/CronTask"
	"time"
)



func (scheduler *TaskScheduler) AddTask(task *CronTask.Task) string {

	if task.RunTime < 100000000000 {
		task.RunTime = task.RunTime * int64(time.Second)
	}
	if task.RunTime < time.Now().UnixNano() {
		task.RunTime = time.Now().UnixNano() + int64(time.Second)
	}

	if task.ID == "" {
		task.ID = uuid.New().String()
	}

	return scheduler.addTask(task)

}

func (scheduler *TaskScheduler) addTask(task CronTask.TaskInterface) string {
	if scheduler.lock {
		scheduler.swap = append(scheduler.swap, task)
	} else {
		scheduler.tasks = append(scheduler.tasks, task)
		scheduler.add <-task
	}

	return task.GetUuid()

}


func (scheduler *TaskScheduler) Lock() {
	scheduler.lock = true
}

func (scheduler *TaskScheduler) UnLock() {
	scheduler.lock = false
	if len(scheduler.swap) > 0 {
		for _, task := range scheduler.swap {
			scheduler.tasks = append(scheduler.tasks, task)
		}
		scheduler.swap = make([]CronTask.TaskInterface, 0)
	}
}