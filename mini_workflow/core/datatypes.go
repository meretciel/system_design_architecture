package core

import (
	"fmt"
	"reflect"
)

type State string

const (
	StateNew         State = "New"
	StateUnknown     State = "Unknown"
	StatePending     State = "Pending"
	StateStarted     State = "Started"
	StateSucceeded   State = "Succeeded"
	StateFailed      State = "Failed"
	StateCancelled   State = "Cancelled"
	StateTimedOut    State = "TimedOut"
	StateInterrupted State = "Interrupted"
)


type Task interface {
	// Core
	GetState() State
	SetState(s State)
	GetDependencies() map[string]*Task
	GetDependency(name string) *Task
	AddDependency(task *Task)
	AddNamedDependency(name string, task *Task)
	
	// Customization
	SerializeParam() string
	LoadParam(param string) error

	SerializeResult() string
	LoadResult() error

	Run() error
	CheckState() State
}

type TaskCore struct {
	state        State
	result       any
	dependencies map[string]*Task
}


func (t *TaskCore) GetState() State {
	return t.state
}

func (t *TaskCore) SetState(s State) {
	t.state = s
}

func (t *TaskCore) GetDependencies() map[string]*Task {
	return t.dependencies
}

func (t *TaskCore) GetDependency(name string) *Task {
	v, exists := t.dependencies[name]
	if !exists {
		panic(fmt.Sprintf("Dependency %s does not exist", name))
	}
	return v
}

func (t *TaskCore) AddNamedDependency(name string, task *Task) {
	_, exists := t.dependencies[name]
	if exists {
		panic(fmt.Sprintf("Dependency %s already exists.", name))
	}
	t.dependencies[name] = task
}

func (t *TaskCore) AddDependency(task *Task) {
	var typ = reflect.TypeOf(task)
	var name = fmt.Sprintf("%s.%s", typ.PkgPath(), typ.Name())
	t.AddNamedDependency(name, task)
}





type Workflow struct {
	Tasks []Task
	CriticalTasks []*Task
}

func (w *Workflow) GetState() State {
	var state = StateNew

	for _, task := range w.Tasks {
		if task.GetState() != StateNew {
			state = StatePending
			break
		}
	}

	if state == StateNew {
		return StateNew
	}

	// state == StatePending
	var counts = make(map[State]int)

	for _, task := range w.CriticalTasks {
		counts[(*task).GetState()] = counts[(*task).GetState()] + 1	
	}

	if counts[StateSucceeded] == len(w.CriticalTasks) {
		return StateSucceeded
	}

	if counts[StateFailed] > 0 {
		return StateFailed
	}

	if counts[StateInterrupted] > 0 {
		return StateInterrupted
	}

	if counts[StateTimedOut] > 0 {
		return StateTimedOut
	}


	if counts[StateCancelled] > 0 {
		return StateCancelled
	}

	return StatePending

}


type TaskExecution struct {
	RunId string
	DependencyTaskIds map[string]string
	TaskParameter string
	Task *Task
}

type WorkflowExecution struct {
	// Runtime information
	RunId string
	WorkflowName string
	TaskRunIds []string
	CriticalTaskRunIds []string
	// In memory instance
	Workflow *Workflow
}






// func (t *GenerateReport) Run() error {
// 	var taskReadBook = t.GetDependency(reflect.TypeFor[ReadBookTask]()).(ReadBookTask) 
// 	var result = taskReadBook.GetResult()	// because GetResult is not in the interface, it has type now
// }

// func ExampleWorkflow() Workflow {
// 	// User code
// 	var readBookParam = ""
// 	var readBookTask = TaskRegistry.Create(reflect.TypeFor[ReadBook]()).SetParameter(readBookParam)

// 	var genReportParam = ""
// 	var genReportTask = TaskRegistry.Create(reflect.TypeFor[GenerateRepot]()).SetParameter(genReportParam)

// 	// Engine code
// 	var paramStr = ""
// 	var readBookTask = TaskRegistry.CreateFromName("ReadBook")
// 	readBookTask.LoadParam(paramStr)
	
// 	var genReportParamStr = ""
// 	var genReportTask = TaskRegistry.CreateFromName("GenerateReport")
// 	genReportTask.LoadParam(genReportParamStr)
// 	genReportTask.AddNamedDependency("ReadBook", readBookTask)



// }
	


// 	var taskGenerateReport = TaskRegistry.Create("GenerateReport").SetParameter(generateReportParam)
// 	taskGenerateReport.AddDependency(taskReadBook)
// }


// func ExampleDefineTasks() {

// 	var tf = TaskFactory {
// 		"mock-task",
// 		reflect.TypeFor[time.Duration](),
// 		func (arg any) task.MockTask {
// 			var typedArg = arg.(time.Duration)
// 			return task.CreateMockTask(typedArg)
// 		},

// 	}

// }