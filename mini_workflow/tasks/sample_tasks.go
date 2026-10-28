package tasks

import (
	wf "mini_workflow"
	"time"
)

type State = wf.State
type Task = wf.Task

type BasicTask struct {
	name         string
	state        State
	result       any
	dependencies []Task
}

func (t *BasicTask) GetName() string {
	return t.name
}

func (t *BasicTask) GetState() State {
	return t.state
}

func (t *BasicTask) SetState(newState State) {
	t.state = newState
}

func (t *BasicTask) GetResult() any {
	return t.result
}

func (t *BasicTask) GetDependencies() []Task {
	return t.dependencies
}

type MockTask struct {
	name     string
	state    State
	execTime time.Duration
}

func (t *MockTask) GetName() string {
	return t.name
}

func (t *MockTask) GetState() State {
	return t.state
}

func (t *MockTask) SetState(newState State) {
	t.state = newState
}

func (t *MockTask) GetResult() any {
	return nil
}

func (t *MockTask) GetDependencies() []Task {
	return nil
}

func (t *MockTask) Run() error {
	time.Sleep(t.execTime)
	return nil
}


func CreateMockTask(execTime time.Duration) MockTask {
	return MockTask{
		"",
		wf.StateUnknown,
		execTime,
	}
}

// type Sequence struct {
// 	BasicTask
// 	tasks []Task
// }

// func NewSequence(name string, tasks ...Task) Task {
// 	if len(sequence) == 0 {
// 		panic("Invalid sequence input: Length is zero.")
// 	}
// 	return Sequence{
// 		name:         name,
// 		state:        wf.StateUnknown,
// 		result:       nil,
// 		dependencies: nil,
// 		tasks: tasks
// 	}
// }

// func (s *Sequence) Run() error {

// 	for _, task := range s.tasks {
// 		var err = task.Run()
// 		s.result = task.GetResult()

// 		if err ==  nil {
// 			s.state = wf.StateSucceeded
// 		} else {
// 			s.state = wf.StateFailed
// 		}
// 	}
// }
