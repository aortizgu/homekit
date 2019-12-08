package devicemanager

import (
	"encoding/json"
	"homekit/app/calderadevice"
	"homekit/app/devicerules"
	"homekit/app/led"
	"homekit/app/msgbroker"
	"homekit/app/notifier"
	"homekit/app/tempsensor"
	"log"
	"net/http"
	"time"
)

// DeviceManager manages the devices
type DeviceManager struct {
}

const hysteresis float64 = 1.0        // 1 degree of hysteresis
const externalWeatherPeriod = 30 * 60 // half hour

var lastExternalWeatherRequest = 0
var LastExternalWeatherTemp = 0.0

func (c DeviceManager) evaluateState(temp float64, activeNow bool) (bool, bool) {
	active := false
	manual := false
	if rules, err := devicerules.DeviceRules.GetWeekRule(calderadevice.Caldera); err == nil {
		limitTemp := rules.GetFloatTemp()
		mustActivate := false
		if activeNow {
			mustActivate = temp < limitTemp+hysteresis
		} else {
			mustActivate = temp < limitTemp-hysteresis
		}
		if mustActivate {
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
		active, manual = c.evaluateState(sensorTemp, calderadevice.CalderaActive)
	}

	err = calderadevice.SetState(active)
	calderadevice.CheckError(err)
	calderaTemp, err := calderadevice.GetTemp()
	calderadevice.CheckError(err)

	status = !calderadevice.CalderaError && !tempsensor.SensorError

	if calderadevice.CalderaActive != active {
		notifier.NotifyCalderaState(active, manual)
		calderadevice.CalderaActive = active
	}

	if int(time.Now().Unix())-lastExternalWeatherRequest > externalWeatherPeriod {
		resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?units=metric&q=Fuenlabrada,es&appid=62b6faef972916a25c2420b17af38d40")
		defer resp.Body.Close()
		if err == nil {
			lastExternalWeatherRequest = int(time.Now().Unix())
			var info map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&info)
			main := info["main"].(map[string]interface{})
			LastExternalWeatherTemp = main["temp"].(float64)
		} else {
			log.Println("Cannot get external temp", err)
		}

	}
	e := msgbroker.NewEvent(sensorTemp, calderaTemp, LastExternalWeatherTemp, status, active, manual)
	msgbroker.Publish(e)
	led.Update(e)
}

// Refresh the status
func Refresh() {
	d := DeviceManager{}
	d.Run()
}
