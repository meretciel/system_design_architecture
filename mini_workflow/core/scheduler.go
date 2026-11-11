package core

import (
	"fmt"
	"log/slog"
)

type Scheduler struct {
	Mailbox chan any
}

func CreateScheduler() Scheduler {
	return Scheduler{
		Mailbox: make(chan any),
	}
}

func (s *Scheduler) Run() {
	// for {
	// 	select {
	// 	case msg := <-s.Mailbox:
	// 		slog.Info(fmt.Sprintf("Received a message %v", msg),
	// 			COMPONENT, SCHEDULER)
	// 	}
	// }

	for msg := range s.Mailbox {
		slog.Info(fmt.Sprintf("Scheduler received a message %#v", msg), COMPONENT, SCHEDULER)
	}
}
