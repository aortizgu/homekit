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
	MONDAY    = iota
	TUESDAY   = iota
	WEDNESDAY = iota
	THURSDAY  = iota
	FRIDAY    = iota
	SATURDAY  = iota
	SUNDAY    = iota
	WEEK_DAYS = iota
)

type TimeRule struct {
	Hour   int
	Minute int
}

type DayRule struct {
	Enabled bool
	Name    string
	Start   TimeRule
	End     TimeRule
}

type WeekRule struct {
	DeviceName      string
	TemperatureInt  int
	TemperatureFrac int
	Manual          bool
	Days            [WEEK_DAYS]DayRule
}

func (c WeekRule) GetFloatTemp() float64 {
	if c.TemperatureFrac == 0 {
		return float64(c.TemperatureInt)
	} else {
		return float64(c.TemperatureInt) + float64(c.TemperatureFrac)/10
	}
}

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

func getDevices() []string {
	devices := make([]string, 0, len(DeviceRules))
	for device := range DeviceRules {
		devices = append(devices, device)
	}
	return devices
}

func (c Rules) Index() revel.Result {
	devices := getDevices()
	return c.Render(devices)
}

func (c Rules) GetDeviceRules(device string) revel.Result {
	if rules, ok := DeviceRules[device]; ok {
		devices := getDevices()
		return c.Render(device, devices, rules)
	}
	return c.Redirect(routes.Rules.Index())
}

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
		for i := 0; i < WEEK_DAYS; i++ {
			iStr := strconv.Itoa(i)
			d := DayRule{}
			d.Name = weekDaysName[i]
			enabled := c.Params.Get("enabled-" + iStr)
			d.Enabled = enabled == "on"

			s := TimeRule{}
			start := c.Params.Get("start-" + iStr)
			if h, m, err := getHMFromString(start); err == nil {
				s.Hour = h
				s.Minute = m
				d.Start = s
			}

			e := TimeRule{}
			end := c.Params.Get("end-" + iStr)
			if h, m, err := getHMFromString(end); err == nil {
				e.Hour = h
				e.Minute = m
				d.End = e
			}

			rules.Days[i] = d
		}

		DeviceRules[device] = rules
		storeConfig()
		return c.Redirect(routes.Rules.GetDeviceRules(device))
	}
	return c.Redirect(routes.Rules.Index())
}

var (
	DeviceRules  = make(map[string]WeekRule)
	weekDaysName = []string{
		"Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado", "Domingo",
	}
)

const (
	ConfigFilename = "rules.json"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func storeConfig() {
	jsonFile, _ := json.MarshalIndent(DeviceRules, "", " ")
	_ = ioutil.WriteFile(ConfigFilename, jsonFile, 0644)
}

func loadConfig() {
	jsonFile, err := os.Open(ConfigFilename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &DeviceRules)
}

func loadRules() revel.Result {
	if fileExists(ConfigFilename) {
		loadConfig()
		log.Println(DeviceRules)
	} else {
		rule := WeekRule{
			DeviceName: "Calefacción",
		}
		for index := 0; index < WEEK_DAYS; index++ {
			rule.Days[index].Name = weekDaysName[index]
			rule.Days[index].Enabled = false
		}
		DeviceRules["caldera"] = rule
		storeConfig()
	}
	return nil
}

func init() {
	loadRules()
}
