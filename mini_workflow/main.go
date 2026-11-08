package main

import (
	"mini_workflow/tasks"
	"mini_workflow/workflows"
	"mini_workflow/core"
)

func main() {
	var engine = core.CreateWorkflowEngine()

	engine.TaskRegistry.Register(core.CastTaskProvider(tasks.CreateSleepTask))

	var workflowProviders = []func() core.Workflow{
		workflows.Demo,
	}

	for _, wfProvider := range workflowProviders {
		engine.Coordinator.Register(wfProvider)
	}

	engine.Run()
}
