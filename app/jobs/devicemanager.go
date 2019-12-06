package jobs

import (
	"homekit/app/calderadevice"
	"homekit/app/controllers"
	"homekit/app/msgbroker"
	"homekit/app/notifier"
	"homekit/app/tempsensor"
	"time"

	"github.com/revel/modules/jobs/app/jobs"
	"github.com/revel/revel"
)

// DeviceManager manages the devices
type DeviceManager struct {
}

func (c DeviceManager) evaluateState(temp float64) (bool, bool) {
	active := false
	manual := false
	if rules, ok := controllers.DeviceRules[calderadevice.Caldera]; ok {
		limitTemp := rules.GetFloatTemp()
		if temp < limitTemp {
			if rules.Manual {
				active = true
				manual = true
			} else {
				weekday := int(time.Now().Weekday())
				if weekday == 0 {
					weekday = 6
				} else {
					weekday = weekday - 1
				}
				dayRule := rules.Days[weekday]
				if dayRule.Enabled {
					hours, minutes, _ := time.Now().Clock()
					minStart := dayRule.Start.Hour*60 + dayRule.Start.Minute
					minStop := dayRule.End.Hour*60 + dayRule.End.Minute
					minCurrent := hours*60 + minutes
					if minCurrent >= minStart && minCurrent <= minStop {
						active = true
					}
				}
			}
		}
	}
	return active, manual
}

// Run runnable method of DeviceManager
func (c DeviceManager) Run() {
	active := false
	manual := false
	status := false
	sensorTemp, err := tempsensor.GetTemp()
	tempsensor.CheckError(err)
	if err == nil {
		active, manual = c.evaluateState(sensorTemp)
	}

	// calderadevice.setState(active)
	// calderadevice.CheckError(err)
	calderaTemp, err := calderadevice.GetTemp()
	calderadevice.CheckError(err)

	status = !calderadevice.CalderaError && !tempsensor.SensorError

	if calderadevice.CalderaActive != active {
		notifier.NotifyCalderaState(active, manual)
		calderadevice.CalderaActive = active
	}

	msgbroker.Publish(msgbroker.NewEvent(sensorTemp, calderaTemp, status, active, manual))
}

func init() {
	revel.OnAppStart(func() {
		jobs.Now(DeviceManager{})
		jobs.Schedule("@every 20s", DeviceManager{})
	})
}
