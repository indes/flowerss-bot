package task

import "go.uber.org/zap"

var (
	taskList = []Task{}
)

type Task interface {
	Start()
	Stop()
	Name() string
}

func registerTask(task Task) {
	taskList = append(taskList, task)
}

func StartTasks() {
	for _, task := range taskList {
		zap.S().Infof("Start task %v", task.Name())
		task.Start()
	}
}

func StopTasks() {
	for _, task := range taskList {
		zap.S().Infof("Stop task %v", task.Name())
		task.Stop()
	}
}
