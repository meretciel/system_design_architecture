package mini_workflow

import (
	"reflect"
    "time"
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
	GetState() State
	SetState(s State)
	CheckState() State
	Run() error
	
	GetDependencies() []Task
	GetDependency(name string) Task

	AddDependency(task Task)
	AddNamedDependency(name string, task Task)
	
	SerializeParam() string
	LoadParam(param string) error

	SerializeResult() string
	LoadResult() error
}

type TaskCore struct {
	state        State
	result       any
	dependencies []Task
}


type Workflow struct {
	Tasks []Task
	CriticalTasks []Task
	State State
	// Schedule Schedule
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