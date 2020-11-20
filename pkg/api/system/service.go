package system

import (
	"fmt"
	"time"

	"github.com/mythio/battery-notif/pkg/common/util/battery"
	"github.com/mythio/battery-notif/pkg/common/util/notify"
)

// System represents the main service
type System struct {
	battery battery.Service
	notify  notify.Service
}

// InitService initializes and returns an instance of the service
func InitService(batterySvc battery.Service, notifySvc notify.Service) *System {
	return &System{
		battery: batterySvc,
		notify:  notifySvc,
	}
}

// Run runs the service
func (s System) Run() {
	for range time.NewTicker(1 * time.Second).C {
		status, _ := s.battery.GetCurrentBatteryPercentage()
		fmt.Println(status)
		if status.Percentage >= 80 && status.IsCharging {
			s.notify.SendNotification("Disconnect from charging for keep battery healthy")
		} else if status.Percentage <= 20 && !status.IsCharging {
			s.notify.SendNotification("Connect charger for keep battery healthy")
		}
	}
}
