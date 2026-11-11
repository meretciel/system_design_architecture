package core

import (
	"fmt"
	"log/slog"
	"time"
)

type Worker struct {
	WorkerId     string
	Mailbox      chan any
	scheduler    chan<- any
	state        WorkerState
	tickInterval time.Duration
	taskRunId    *string
}

type msgTaskIsDone struct{}

func NewWorker(workerId string, scheduler chan<- any, tickInterval time.Duration) Worker {
	return Worker{
		WorkerId:     workerId,
		Mailbox:      make(chan any),
		scheduler:    scheduler,
		tickInterval: tickInterval,
		state:        WorkerStateIdle,
		taskRunId:    nil,
	}
}

func (w *Worker) Run() {
	slog.Info("Worker starts.", COMPONENT, WORKER, WORKER_ID, w.WorkerId)

	var ticker = time.NewTicker(w.tickInterval)
	defer ticker.Stop()

	for {
		select {
		case msg := <-w.Mailbox:
			w.handleMessage(msg)

		case <-ticker.C:
			w.scheduler <- WorkerReport{
				Timestamp: time.Now().UnixMilli(),
				Id:        w.WorkerId,
				State:     w.state,
				TaskRunId: w.taskRunId,
			}
		}
	}
}

func (w *Worker) handleMessage(msg any) {
	slog.Info(fmt.Sprintf("Worker receives message: %#v", msg), COMPONENT, WORKER, WORKER_ID, w.WorkerId)
	switch v := msg.(type) {
	// Internal messages
	case msgTaskIsDone:
		w.state = WorkerStateIdle
	// External messages
	case ReqRunTask:
		// We assume the scheduler does the book keeping in this implementation
		// Worker does not queue requests. This could be problematic if it receives
		// tasks when it's is already running.
		switch w.state {
		case WorkerStateIdle:
			w.state = WorkerStateRunning
			go func() {
				v.Task.Run()
				w.Mailbox <- msgTaskIsDone{}
			}()
			w.scheduler <- ResponseRunTask{
				ReqId:         v.ReqId,
				ResponseState: ResponseStateAck,
			}

		case WorkerStateRunning:
			w.scheduler <- ResponseRunTask{
				ReqId:         v.ReqId,
				ResponseState: ResponseStateFailed,
			}
		}
	}
}
