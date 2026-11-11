package core

import (
	"fmt"
	"reflect"
	"runtime"

)

func CastTaskProvider[T Task](provider func() T) func() Task {
	return func() Task {
		return Task(provider())
	}
}


// This is the unique identifier of a Task.
// It's the package + task struct name
func GetTaskSpecIdentifier(t Task) string {
	var typ = reflect.TypeOf(t)
	return fmt.Sprintf("%s.%s", typ.PkgPath(), typ.Name())
}


func GetWorkflowSpecIdentifier(w func() Workflow) string {
	var p = reflect.ValueOf(w).Pointer()
	return runtime.FuncForPC(p).Name()
}

