package controllers

import (
	"errors"
	"homekit/app/devicemanager"
	"homekit/app/devicerules"
	"homekit/app/routes"
	"log"
	"strconv"
	"strings"

	"github.com/revel/revel"
)

func getHMFromString(s string) (int, int, error) {
	log.Println("getHMFromString", s)
	hm := strings.Split(s, ":")
	if len(hm) != 2 {
		return 0, 0, errors.New("invalid size")
	}
	h, err := strconv.Atoi(hm[0])
	if err != nil {
		return 0, 0, err
	}
	m, err := strconv.Atoi(hm[1])
	if err != nil {
		return 0, 0, err
	}
	return h, m, nil
}

// Rules web page to change configuration for devices
type Rules struct {
	Application
}

func (c Rules) checkUser() revel.Result {
	log.Println("checkuser")
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

// Index shows index page of rules
func (c Rules) Index() revel.Result {
	devices := devicerules.DeviceRules.GetDevices()
	return c.Render(devices)
}

// GetDeviceRules shows rules for an input device
func (c Rules) GetDeviceRules(device string) revel.Result {
	if rules, err := devicerules.DeviceRules.GetWeekRule(device); err == nil {
		devices := devicerules.DeviceRules.GetDevices()
		return c.Render(device, devices, rules)
	}
	return c.Redirect(routes.Rules.Index())
}

// SetDeviceRules set rules for an input device
func (c Rules) SetDeviceRules(device string) revel.Result {
	if rules, err := devicerules.DeviceRules.GetWeekRule(device); err == nil {
		c.Log.Info("Seting device rules  ", c.Params)
		tInt, err := strconv.Atoi(c.Params.Get("temp-int"))
		if err != nil {
			c.Log.Error("Error in conversion of ", c.Params.Get("temp-int"))
		} else {
			rules.TemperatureInt = tInt
		}
		tFrac, err := strconv.Atoi(c.Params.Get("temp-frac"))
		if err != nil {
			c.Log.Error("Error in conversion of ", c.Params.Get("temp-int"))
		} else {
			rules.TemperatureFrac = tFrac
		}
		manual := c.Params.Get("manual")
		rules.Manual = manual == "on"
		if !rules.Manual {
			for i := 0; i < devicerules.WeekDays; i++ {
				iStr := strconv.Itoa(i)
				d := devicerules.DayRule{}
				d.Name = devicerules.WeekDaysName[i]
				enabled := c.Params.Get("enabled-" + iStr)
				d.Enabled = enabled == "on"

				start := c.Params.Get("start-" + iStr)
				s := devicerules.TimeRule{0, start}
				if h, m, err := getHMFromString(start); err == nil {
					s.Mins = h*60 + m
					d.Start = s
				}

				end := c.Params.Get("end-" + iStr)
				e := devicerules.TimeRule{0, end}
				if h, m, err := getHMFromString(end); err == nil {
					e.Mins = h*60 + m
					d.End = e
				}

				rules.Days[i] = d
			}
		}
		devicerules.DeviceRules.SetWeekRule(device, rules)
		devicerules.DeviceRules.StoreRules()
		go devicemanager.Refresh()
		return c.Redirect(routes.Rules.GetDeviceRules(device))
	}
	return c.Redirect(routes.Rules.Index())
}

func init() {
	devicerules.DeviceRules.LoadRules()
}
