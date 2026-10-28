package mini_workflow


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


func (engine *WorkflowEngine) Run() {
	
}