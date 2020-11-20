package api

import (
	"github.com/mythio/battery-notif/pkg/api/system"
	"github.com/mythio/battery-notif/pkg/common/util/battery"
	"github.com/mythio/battery-notif/pkg/common/util/notify"
)

// Start starts the service
func Start() error {
	batterySvc := battery.New()
	notifySvc := notify.New()

	sys := system.InitService(batterySvc, notifySvc)
	sys.Run()

	return nil
}
