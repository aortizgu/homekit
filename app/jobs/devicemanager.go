package jobs

import (
	"homekit/app/devicemanager"

	"github.com/revel/modules/jobs/app/jobs"
	"github.com/revel/revel"
)

func init() {
	revel.OnAppStart(func() {
		jobs.Now(devicemanager.DeviceManager{})
		jobs.Schedule("@every 20s", devicemanager.DeviceManager{})
	})
}
