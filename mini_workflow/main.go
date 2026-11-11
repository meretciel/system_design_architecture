package main

import (
	"mini_workflow/tasks"
	"mini_workflow/workflows"
	"mini_workflow/core"
	"log/slog"
	"time"
)

func main() {


	var engine = core.NewWorkflowEngine(core.WorkflowEngineConfig{
		NumWorkers: 2,
		WorkerTickInterval: 5 * time.Second,
		WorkflowPollInterval: 10 * time.Second,
	})

	engine.TaskRegistry.Register(core.CastTaskProvider(tasks.CreateSleepTask))

	var workflowProviders = []func() core.Workflow{
		workflows.Demo,
	}

	for _, wfProvider := range workflowProviders {
		engine.Coordinator.Register(wfProvider)
	}

	engine.Run()

	slog.Info("Program ends")
}
