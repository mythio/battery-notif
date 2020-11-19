package notify

import (
	"context"
	"os/exec"
	"time"
)

// Service provides notification sender implementation
type Service struct{}

// New creates a new service instance
func New() Service {
	return Service{}
}

// SendNotification sends a notification
func (s Service) SendNotification(text string) error {
	ch := make(chan error)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		notify := exec.CommandContext(ctx, "notify-send", text)
		if err := notify.Run(); err != nil {
			ch <- err
		}

		ch <- nil
	}()

	err := <-ch
	return err
}
