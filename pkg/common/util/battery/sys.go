package battery

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Service provides battery percentage extraction implementation
type Service struct{}

type result struct {
	result uint16
	err    error
}

// New creates a new service instance
func New() Service {
	return Service{}
}

// GetCurrentBatteryPercentage returns current battery percentage
func (s Service) GetCurrentBatteryPercentage() (uint16, error) {
	ch := make(chan result, 1)

	go func() {
		var err error
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		upower := exec.CommandContext(ctx, "upower", "-i", "/org/freedesktop/UPower/devices/battery_BAT0")
		grep := exec.CommandContext(ctx, "grep", "percentage")

		grep.Stdin, err = upower.StdoutPipe()
		if err != nil {
			ch <- result{0, errors.New("grep pipe error")}
		}

		out := &bytes.Buffer{}
		grep.Stdout = out

		if err = grep.Start(); err != nil {
			ch <- result{0, errors.New("grep exec error")}
			return
		}
		if err = upower.Run(); err != nil {
			ch <- result{0, errors.New("grep exec error")}
			return
		}

		res := out.String()

		percentageStr := strings.Trim(res[len(res)-4:len(res)-1], " ")
		percentage, err := strconv.ParseInt(percentageStr, 10, 16)
		if err != nil {
			ch <- result{0, errors.New("grep exec error")}
		}

		ch <- result{uint16(percentage), nil}
	}()

	res := <-ch
	return res.result, nil
}
