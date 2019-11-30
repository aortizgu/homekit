package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"homekit/app/routes"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/revel/revel"
)

const (
	configFilename string = "rules.json"
	weekDays       int    = 7
)

type timeRule struct {
	Hour   int
	Minute int
}

type dayRule struct {
	Enabled bool
	Name    string
	Start   timeRule
	End     timeRule
}

type weekRule struct {
	DeviceName      string
	TemperatureInt  int
	TemperatureFrac int
	Manual          bool
	Days            [weekDays]dayRule
}

func (c weekRule) GetFloatTemp() float64 {
	if c.TemperatureFrac == 0 {
		return float64(c.TemperatureInt)
	}
	return float64(c.TemperatureInt) + float64(c.TemperatureFrac)/10
}

// DeviceRulesStorage storage of rules inxed by device name
type DeviceRulesStorage map[string]weekRule

func (c DeviceRulesStorage) getDevices() []string {
	devices := make([]string, 0, len(c))
	for device := range c {
		devices = append(devices, device)
	}
	return devices
}

func (c DeviceRulesStorage) storeConfig() {
	jsonFile, _ := json.MarshalIndent(DeviceRules, "", " ")
	_ = ioutil.WriteFile(configFilename, jsonFile, 0644)
}

func (c DeviceRulesStorage) loadConfig() {
	jsonFile, err := os.Open(configFilename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &DeviceRules)
}

func (c DeviceRulesStorage) loadRules() revel.Result {
	if fileExists(configFilename) {
		c.loadConfig()
		log.Println(DeviceRules)
	} else {
		rule := weekRule{
			DeviceName: "Calefacción",
		}
		for index := 0; index < weekDays; index++ {
			rule.Days[index].Name = weekDaysName[index]
			rule.Days[index].Enabled = false
		}
		DeviceRules["caldera"] = rule
		c.storeConfig()
	}
	return nil
}

var (
	// DeviceRules configuration of rules of devices
	DeviceRules  = make(DeviceRulesStorage)
	weekDaysName = []string{
		"Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado", "Domingo",
	}
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
	devices := DeviceRules.getDevices()
	return c.Render(devices)
}

// GetDeviceRules shows rules for an input device
func (c Rules) GetDeviceRules(device string) revel.Result {
	if rules, ok := DeviceRules[device]; ok {
		devices := DeviceRules.getDevices()
		return c.Render(device, devices, rules)
	}
	return c.Redirect(routes.Rules.Index())
}

// SetDeviceRules set rules for an input device
func (c Rules) SetDeviceRules(device string) revel.Result {
	if rules, ok := DeviceRules[device]; ok {
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
		for i := 0; i < weekDays; i++ {
			iStr := strconv.Itoa(i)
			d := dayRule{}
			d.Name = weekDaysName[i]
			enabled := c.Params.Get("enabled-" + iStr)
			d.Enabled = enabled == "on"

			s := timeRule{}
			start := c.Params.Get("start-" + iStr)
			if h, m, err := getHMFromString(start); err == nil {
				s.Hour = h
				s.Minute = m
				d.Start = s
			}

			e := timeRule{}
			end := c.Params.Get("end-" + iStr)
			if h, m, err := getHMFromString(end); err == nil {
				e.Hour = h
				e.Minute = m
				d.End = e
			}

			rules.Days[i] = d
		}

		DeviceRules[device] = rules
		DeviceRules.storeConfig()
		return c.Redirect(routes.Rules.GetDeviceRules(device))
	}
	return c.Redirect(routes.Rules.Index())
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func init() {
	DeviceRules.loadRules()
}
