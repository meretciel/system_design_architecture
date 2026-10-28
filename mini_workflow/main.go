package mini_workflow

import (
	"workflows"
	"tasks"
)

func main() {
	var engine = CreateWorkflowEngine()

	engine.TaskRegistry.Register(tasks.CreateSleepTask)


	
	var workflows = []func() Workflow {
		workflows.Demo,
		workflows.Demo2,
	}

	for _, wfProvider := range workflows {
		engine.Coordinator.Register(wfProvider)
	}

	engine.Run()
}
