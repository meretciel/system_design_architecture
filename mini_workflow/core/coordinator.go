package core

import (
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type Coordinator struct {
	taskRegistry         TaskRegistry
	workflowProviders    map[string]func() Workflow   // map name to workflow tempalte/provider
	workflowExecutions   map[string]WorkflowExecution // map name to workflow execution
	taskExecutions       map[Task]*TaskExecution
	// taskRunIdToExecution map[string]*TaskExecution
}

func CreateCoordinator() Coordinator {
	return Coordinator{
		taskRegistry:         CreateTaskRegistry(),
		workflowProviders:    make(map[string]func() Workflow),
		workflowExecutions:   make(map[string]WorkflowExecution),
		taskExecutions:       make(map[Task]*TaskExecution),
		// taskRunIdToExecution: make(map[string]*TaskExecution),
	}
}

func (co *Coordinator) Register(workflowProvider func() Workflow) {
	var name = GetWorkflowSpecIdentifier(workflowProvider)
	co.workflowProviders[name] = workflowProvider
}

func (co *Coordinator) loadFromPersistenceStore() {

}

func (co *Coordinator) Run() {
	slog.Info("Coordinator starts to run")
	// In the current implemenation, we don't do scheudling but
	// we still need to load from the persistence representation layer
	co.loadFromPersistenceStore()

	// Check the execution state of workflow
	for name, workflowProvider := range co.workflowProviders {
		_, exists := co.workflowExecutions[name]
		if !exists {
			slog.Info("No workflow execution found. Will create one if needed.", WORKFLOW, name)

			// Create a new execution instance in memory
			var workflow = workflowProvider()
			var taskExecutions = []*TaskExecution{}

			// Initialize the task executions records
			for i := range workflow.Tasks {
				var task = workflow.Tasks[i]
				var taskRunId = uuid.New().String()
				var taskExecution = TaskExecution{
					RunId:                taskRunId,
					TaskParameter:        task.SerializeParam(),
					Task:                 task,
					DependencyTaskRunIds: make(map[string]string),
				}

				slog.Info(fmt.Sprintf("Create task execution for task=%p, taskRunId=%s", task, taskRunId))

				taskExecutions = append(taskExecutions, &taskExecution)
				co.taskExecutions[task] = &taskExecution
				// co.taskRunIdToExecution[taskRunId] = co.taskExecutions[task]
			}

			// Fill in dependencies
			for i := range workflow.Tasks {
				var task = workflow.Tasks[i]
				var taskExecution = co.taskExecutions[task]
				for name, dep := range task.GetDependencies() {
					slog.Info(fmt.Sprintf("name=%s, dep=%p", name, dep))
					(*taskExecution).DependencyTaskRunIds[name] = co.taskExecutions[dep].RunId
				}
			}

			var criticalTaskExecutions = []*TaskExecution{}

			for _, task := range workflow.CriticalTasks {
				var execution = co.taskExecutions[task]
				criticalTaskExecutions = append(criticalTaskExecutions, execution)
			}

			var workflowExecution = WorkflowExecution{
				RunId:                  uuid.New().String(),
				WorkflowName:           name,
				Workflow:               &workflow,
				TaskExecutions:         taskExecutions,
				CriticalTaskExecutions: criticalTaskExecutions,
			}

			co.workflowExecutions[name] = workflowExecution
			slog.Info("Created a new workflow execution.", WORKFLOW, name)
		}
	}
}
