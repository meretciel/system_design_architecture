package mini_workflow

import (
	"reflect"
	"testing"
)

type registryTestTask struct {
	Task
	id int
}

func TestTaskRegistryCreate(t *testing.T) {
	var registry = CreateTaskRegistry()
	count := 0
	registry.Register(func() Task {
		count++
		return &registryTestTask{id: count}
	})

	first := registry.Create(reflect.TypeOf(&registryTestTask{})).(*registryTestTask)
	second := registry.CreateByName(reflect.TypeOf(&registryTestTask{}).String()).(*registryTestTask)
	if first == second || first.id != 1 || second.id != 2 {
		t.Fatalf("expected separate tasks from the creator, got %#v and %#v", first, second)
	}
}

func TestTaskRegistryRejectsInvalidRequests(t *testing.T) {
	tests := []struct {
		name string
		run  func(*TaskRegistry)
	}{
		{"nil creator", func(r *TaskRegistry) { r.Register(nil) }},
		{"nil type", func(r *TaskRegistry) { r.Create(nil) }},
		{"unknown name", func(r *TaskRegistry) { r.CreateByName("missing") }},
		{"duplicate registration", func(r *TaskRegistry) {
			creator := func() Task { return &registryTestTask{} }
			r.Register(creator)
			r.Register(creator)
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Error("expected a panic")
				}
			}()
			tt.run(&TaskRegistry{creators: make(map[string]taskProvider)})
		})
	}
}
