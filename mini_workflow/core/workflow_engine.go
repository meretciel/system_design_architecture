package core

import (
	"fmt"
	"log/slog"
	"time"
)


type WorkflowEngine struct {
	TaskRegistry TaskRegistry
	Coordinator Coordinator
	Scheduler Scheduler
}


func CreateWorkflowEngine() WorkflowEngine {
	return WorkflowEngine {
		TaskRegistry: CreateTaskRegistry(),
		Coordinator: CreateCoordinator(),
		Scheduler: CreateScheduler(),
	}
}


// Validate engine specification.
//   - Verify tasks used in workflows are registered
func (engine *WorkflowEngine) ValidateWorkflow(w *Workflow) error {
	for i := range w.Tasks {
		var name = GetTaskSpecIdentifier(w.Tasks[i])
		if !engine.isTaskRegistered(name) {
			return fmt.Errorf("Task %s is not registered in the workflow engine.", name)
		}

	}
	return nil
}


func (engine *WorkflowEngine) isTaskRegistered(name string) bool {
	return engine.TaskRegistry.IsTaskRegistered(name)
}


func (engine *WorkflowEngine) Run() {
	slog.Info("Workflow engine starts to run")
	go engine.Coordinator.Run()
	time.Sleep(time.Second * 5)
}