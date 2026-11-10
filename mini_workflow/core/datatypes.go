package core

import (
	"fmt"
	"reflect"
	"sync"
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
	GetState() State	// This needs to be concurrency safe
	SetState(s State)   // This needs to be concurrency safe
	GetDependencies() map[string]Task
	GetDependency(name string) Task
	AddDependency(task Task)
	AddNamedDependency(name string, task Task)

	// Customization
	SerializeParam() string
	LoadParam(param string) error

	SerializeResult() string
	LoadResult() error

	Run() error
	CheckState() State // This needs to be concurrency safe.
}

// TaskCore provides the core implemenation of the Task interface.
// In interesting implication of using pointer receiver is that
// Only pointers can be passed in the place of Task.
// This has consequence of user code though..
type TaskCore struct {
	state        State
	result       any
	dependencies map[string]Task
	mu sync.Mutex
}

func (t *TaskCore) Init() {
	t.state = StateNew
	t.result = nil
	t.dependencies = make(map[string]Task)
}

func (t *TaskCore) GetState() State {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.state
}

func (t *TaskCore) SetState(s State) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.state = s
}

func (t *TaskCore) GetDependencies() map[string]Task {
	return t.dependencies
}

func (t *TaskCore) GetDependency(name string) Task {
	v, exists := t.dependencies[name]
	if !exists {
		panic(fmt.Sprintf("Dependency %s does not exist", name))
	}
	return v
}

func (t *TaskCore) AddNamedDependency(name string, task Task) {
	_, exists := t.dependencies[name]
	if exists {
		panic(fmt.Sprintf("Dependency %s already exists.", name))
	}
	t.dependencies[name] = task
}

func (t *TaskCore) AddDependency(task Task) {
	var typ = reflect.TypeOf(task)
	var name = fmt.Sprintf("%s.%s", typ.PkgPath(), typ.Name())
	t.AddNamedDependency(name, task)
}

type Workflow struct {
	Tasks         []Task
	CriticalTasks []Task
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
		counts[task.GetState()] = counts[task.GetState()] + 1
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
	RunId                string
	DependencyTaskRunIds map[string]string
	TaskParameter        string
	Task                 Task
}

type WorkflowExecution struct {
	// Runtime information
	RunId              string
	WorkflowName       string
	TaskExecutions         []*TaskExecution
	CriticalTaskExecutions []*TaskExecution
	// In memory instance
	Workflow *Workflow
}


func (w *WorkflowExecution) GetState() State {
	return w.Workflow.GetState()
}