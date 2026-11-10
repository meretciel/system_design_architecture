package tasks

import (
	core "mini_workflow/core"
	"time"
)

type SleepTask struct {
	core.TaskCore
	duration time.Duration
}

// We need to return a pointer because it's the pointer receiver
// that implements the Task interface. SleepTask struct itself
// does not implemnet the Task interface.
func CreateSleepTask() *SleepTask {
	var t = SleepTask{}
	t.Init()
	return &t
}


func (t *SleepTask) SerializeParam() string {
	return ""
}

func (t *SleepTask) LoadParam(param string) error {
	return nil
}

func (t *SleepTask) SerializeResult() string {
	return ""
}

func (t *SleepTask) LoadResult() error {
	return nil
}

func (t *SleepTask) Run() error {
	time.Sleep(t.duration)
	t.SetState(core.StateSucceeded)
	return nil
}


func (t *SleepTask) CheckState() core.State {
	return t.GetState()
}


func (t *SleepTask) SetDuration(duration time.Duration) {
	t.duration = duration
}

