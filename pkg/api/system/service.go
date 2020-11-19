package system

import (
	"github.com/mythio/battery-notif/pkg/common/util/battery"
	"github.com/mythio/battery-notif/pkg/common/util/notify"
)

type System struct {
	Battery battery.Service
	Notify  notify.Service
}

// InitService initializes and returns an instance of the service
func InitService(batterySvc battery.Service, notifySvc notify.Service) *System {
	return &System{
		Battery: batterySvc,
		Notify:  notifySvc,
	}
}
