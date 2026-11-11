package core

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)


type WorkflowEngine struct {
	TaskRegistry TaskRegistry
	Coordinator Coordinator
	Scheduler Scheduler
	config WorkflowEngineConfig
}

type WorkflowEngineConfig struct {
	NumWorkers int
	WorkerTickInterval time.Duration
	WorkflowPollInterval time.Duration
}

func NewWorkflowEngine(config WorkflowEngineConfig) WorkflowEngine {

	var scheduler = CreateScheduler()

	return WorkflowEngine {
		TaskRegistry: CreateTaskRegistry(),
		Coordinator: NewCoordinator(config.WorkflowPollInterval, scheduler.Mailbox),
		Scheduler: scheduler,
		config: config,
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
	go engine.Scheduler.Run()


	var workerInfoList = []WorkerInfo{}
	var workers = []Worker{}

	for range engine.config.NumWorkers {
		var worker = NewWorker(uuid.NewString(), engine.Scheduler.Mailbox, engine.config.WorkerTickInterval)
		workers = append(workers, worker)
		workerInfoList = append(workerInfoList, WorkerInfo{Id: worker.WorkerId, Mailbox: worker.Mailbox})
	}

	for i := range engine.config.NumWorkers {
		go workers[i].Run()
	}

	engine.Scheduler.Mailbox <- MsgWorkerDiscovery {
		WorkerInfoList: workerInfoList,
	}

	time.Sleep(time.Second * 15)
}