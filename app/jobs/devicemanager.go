package jobs

import (
	"homekit/app/controllers"
	"homekit/app/models"
	"homekit/app/msgbroker"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/revel/modules/jobs/app/jobs"
	gorp "github.com/revel/modules/orm/gorp/app"
	"github.com/revel/revel"
)

const (
	deviceManagerStr string = "DeviceManager"
	caldera          string = "caldera"
	domain           string = ".local"
	uptimeURI        string = "/uptime"
	tempURI          string = "/temp"
	relayStateURI    string = "/relay?state="
)

// DeviceManager manages the devices
type DeviceManager struct {
}

var (
	calderaActive bool = false
	calderaError  bool = false
)

func (c DeviceManager) refresh() (string, string, error) {
	uptime := ""
	temp := ""
	resp, err := http.Get("http://" + caldera + domain + uptimeURI)
	if err != nil {
		return uptime, temp, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	uptime = string(body)
	resp, err = http.Get("http://" + caldera + domain + tempURI)
	if err != nil {
		return uptime, temp, err
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	temp = string(body)
	return uptime, temp, err
}

func (c DeviceManager) evaluateState(temp float64) (bool, bool) {
	active := false
	manual := false
	if rules, ok := controllers.DeviceRules[caldera]; ok {
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

func (c DeviceManager) setState(state bool) error {
	active := "off"
	if state {
		active = "on"
	}
	resp, err := http.Get("http://" + caldera + domain + relayStateURI + active)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// Run runnable method of DeviceManager
func (c DeviceManager) Run() {
	uptime, temp, err := c.refresh()
	active := false
	manual := false
	status := "error"
	if err == nil {
		if calderaError {
			controllers.SendMail("Comunicación establecida", "Se ha establecido comunicación con el dispositivo caldera")
		}
		calderaError = false
		tempF, err := strconv.ParseFloat(temp, 64)
		if err == nil {
			status = "ok"
			active, manual = c.evaluateState(tempF)
			//store temp
			meassurement := models.Meassurement{
				Device: caldera,
				Val:    tempF,
				Time:   time.Now().Unix()}
			if err := gorp.Db.Map.Insert(&meassurement); err != nil {
				panic(err)
			}
		}
		err = nil //c.setState(active)
		if err != nil {
			panic(err)
		}
	} else {
		if !calderaError {
			controllers.SendMail("Error de comunicación", "No hay comunicación con el dispositivo caldera")
		}
		calderaError = true
	}
	if calderaActive != active {
		controllers.NotifyCalderaState(active, manual)
		calderaActive = active
	}
	msgbroker.Say(deviceManagerStr, caldera, status, uptime, temp, active)
}

func init() {
	revel.OnAppStart(func() {
		jobs.Now(DeviceManager{})
		jobs.Schedule("@every 10s", DeviceManager{})
	})
}
