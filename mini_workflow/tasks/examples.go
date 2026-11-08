package tasks

import (
	core "mini_workflow/core"
	"sync"
	"time"
)

type SleepTask struct {
	core.TaskCore
	duration time.Duration
	mu sync.Mutex
}

func CreateSleepTask() *SleepTask {
	return &SleepTask{}
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

	t.mu.Lock()
	defer t.mu.Unlock()
	t.SetState(core.StateSucceeded)
	return nil
}


func (t *SleepTask) CheckState() core.State {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.GetState()
}


func (t *SleepTask) SetDuration(duration time.Duration) {
	t.duration = duration
}

