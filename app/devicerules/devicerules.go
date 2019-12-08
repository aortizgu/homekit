package devicerules

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/revel/revel"
)

const (
	WeekDays       int    = 7
	configFilename string = "rules.json"
)

var WeekDaysName = []string{
	"Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado", "Domingo",
}

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
	Days            [WeekDays]DayRule
}

func (c WeekRule) GetFloatTemp() float64 {
	if c.TemperatureFrac == 0 {
		return float64(c.TemperatureInt)
	}
	return float64(c.TemperatureInt) + float64(c.TemperatureFrac)/10
}

// DeviceRulesDevicesMap storage
type DeviceRulesDevicesMap map[string]WeekRule

// DeviceRulesStorage storage of rules inxed by device name
type DeviceRulesStorage struct {
	rules DeviceRulesDevicesMap
}

// NewDeviceRulesStorage new object
func NewDeviceRulesStorage() DeviceRulesStorage {
	ret := DeviceRulesStorage{}
	ret.rules = make(DeviceRulesDevicesMap)
	return ret
}

var (
	// DeviceRules configuration of rules of devices
	DeviceRules    = NewDeviceRulesStorage()
	errFoundDevice = errors.New("failed get device in rules")
)

func (c DeviceRulesStorage) GetWeekRule(device string) (WeekRule, error) {
	if weekRule, ok := DeviceRules.rules[device]; ok {
		return weekRule, nil
	}
	return WeekRule{}, errFoundDevice
}

func (c DeviceRulesStorage) SetWeekRule(device string, w WeekRule) {
	DeviceRules.rules[device] = w
}

func (c DeviceRulesStorage) GetDevices() []string {
	devices := make([]string, 0, len(c.rules))
	for device := range c.rules {
		devices = append(devices, device)
	}
	return devices
}

func (c DeviceRulesStorage) StoreRules() {
	jsonFile, _ := json.MarshalIndent(c.rules, "", " ")
	_ = ioutil.WriteFile(configFilename, jsonFile, 0644)
}

func (c DeviceRulesStorage) LoadRules() revel.Result {
	if fileExists(configFilename) {
		c.loadConfig()
		log.Println(c.rules)
	} else {
		rule := WeekRule{
			DeviceName: "Calefacción",
		}
		for index := 0; index < WeekDays; index++ {
			rule.Days[index].Name = WeekDaysName[index]
			rule.Days[index].Enabled = false
		}
		c.rules["caldera"] = rule
		c.StoreRules()
	}
	return nil
}

func (c DeviceRulesStorage) loadConfig() {
	jsonFile, err := os.Open(configFilename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &c.rules)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
