package mini_workflow

import "reflect"

type Coordinator struct {
	taskRegistry TaskRegistry
}


func (co *Coordinator) Register(workflowProvider func() Workflow) {

}

func (co *Coordinator) Run() {

}

