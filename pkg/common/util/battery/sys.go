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

// Status holds the battery percentage and state
type Status struct {
	IsCharging bool
	Percentage uint16
}

type result struct {
	result Status
	err    error
}

// New creates a new service instance
func New() Service {
	return Service{}
}

// GetCurrentBatteryPercentage returns current battery percentage
func (s Service) GetCurrentBatteryPercentage() (Status, error) {
	ch := make(chan result, 1)

	go func() {
		var err error
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		upower := exec.CommandContext(ctx, "upower", "-i", "/org/freedesktop/UPower/devices/battery_BAT0")
		grep := exec.CommandContext(ctx, "grep", "-e", "percentage", "-e", "state")

		grep.Stdin, err = upower.StdoutPipe()

		if err != nil {
			ch <- result{Status{}, errors.New("grep pipe error")}
		}

		outPercentage := &bytes.Buffer{}
		grep.Stdout = outPercentage

		if err = grep.Start(); err != nil {
			ch <- result{Status{}, errors.New("grep exec error")}
			return
		}
		if err = upower.Run(); err != nil {
			ch <- result{Status{}, errors.New("grep exec error")}
			return
		}

		res := strings.Split(outPercentage.String(), "\n")

		/**
		 * res[0] holds state
		 * res[1] holds percentage
		 */
		percentageStr := strings.Trim(res[1][len(res[1])-4:len(res[1])-1], " ")
		percentage, err := strconv.ParseInt(percentageStr, 10, 16)
		if err != nil {
			ch <- result{Status{}, errors.New("grep exec error")}
		}

		state := strings.Trim(res[0][len(res[0])-14:], " ")

		ch <- result{Status{Percentage: uint16(percentage), IsCharging: state == "charging"}, nil}
	}()

	res := <-ch
	return res.result, nil
}
