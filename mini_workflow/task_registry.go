package mini_workflow

import (
	"fmt"
	"reflect"
)

type taskProvider = func() Task

type TaskRegistry struct {
	creators map[string]taskProvider
}

func CreateTaskRegistry() TaskRegistry {
	return TaskRegistry {
		creators: make(map[string]taskProvider),
	}
}

func (r *TaskRegistry) Register(taskCreator taskProvider) {
	if taskCreator == nil {
		panic("cannot register a nil task creator")
	}


	var name = reflect.TypeOf(taskCreator).Out(0).String()
	_, exist := r.creators[name]

	if exist {
		panic(fmt.Sprintf("Task %s is already registered", name))
	}
	r.creators[name] = taskCreator
}

func (tr *TaskRegistry) Create(taskType reflect.Type) Task {
	if taskType == nil {
		panic("cannot create a task with a nil type")
	}

	return tr.CreateByName(taskType.String())
}

func (tr *TaskRegistry) CreateByName(name string) Task {
	creator, exists := tr.creators[name]
	if !exists {
		panic(fmt.Sprintf("Task %s is not registered", name))
	}

	return creator()
}
