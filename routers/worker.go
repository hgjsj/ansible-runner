package routers

import "ansible-runner/tasks"

func NewWorker() {
	tasks.NewAsyncServer().RegisterHandlers()
	tasks.NewAsyncServer().LaunchServer()
}
