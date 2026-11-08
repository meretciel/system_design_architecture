package core

// import (
// 	"reflect"
// 	"runtime"
// )

import (
	"log/slog"
	"github.com/google/uuid"
)

type Coordinator struct {
	taskRegistry TaskRegistry
	workflowProviders map[string]func() Workflow
	workflowExecutions map[string]WorkflowExecution
}


func CreateCoordinator() Coordinator {
	return Coordinator{
		taskRegistry: CreateTaskRegistry(),
		workflowProviders: make(map[string]func() Workflow),
		workflowExecutions: make(map[string]WorkflowExecution),
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
			// Create a new execution instance in memory
			var workflow = workflowProvider()
			var taskExecutions = make(map[*Task]TaskExecution)
			
			// Initialize the task executions records
			for i := range workflow.Tasks {
				var task = &(workflow.Tasks[i])
				taskExecutions[task] = TaskExecution {
					RunId: uuid.New().String(),
					TaskParameter: (*task).SerializeParam(),
					Task: task,
				}
			}

			// Fill in dependencies
			for i := range workflow.Tasks {
				var task = &(workflow.Tasks[i])
				var taskExecution = taskExecutions[task]
				for name, dep := range (*task).GetDependencies() {
					taskExecution.DependencyTaskIds[name] = taskExecutions[dep].RunId
				}
			}

			var taskRunIds = []string{}

			for _, taskExecution := range taskExecutions {
				taskRunIds = append(taskRunIds, taskExecution.RunId)
			}

			var criticalTaskRunIds = []string{}
			for _, task := range workflow.CriticalTasks {
				criticalTaskRunIds = append(criticalTaskRunIds, taskExecutions[task].RunId)
			}

			var workflowExecution = WorkflowExecution {
				RunId: uuid.New().String(),
				WorkflowName: GetWorkflowSpecIdentifier(workflowProvider),
				Workflow: &workflow,
				TaskRunIds: taskRunIds,
				CriticalTaskRunIds: criticalTaskRunIds,
			}

			co.workflowExecutions[workflowExecution.WorkflowName] = workflowExecution
		}
	}
}




