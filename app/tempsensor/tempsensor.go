package tempsensor

// device ds18b20
import (
	"errors"
	"homekit/app/notifier"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	// SensorError status of error of the device
	SensorError bool = false
)

var errReadSensor = errors.New("failed to read sensor temperature")

// CheckError check if error state and notify
func CheckError(err error) {
	if err == nil {
		if SensorError {
			SensorError = false
			log.Println("calderadevice:: caldera Error ", SensorError)
			if err := notifier.SendMail("Sensor funcionando "+time.Now().Format("15:04:05"), "Se ha establecido comunicación con el sensor de temperatura"); err != nil {
				log.Println("Cannot send mail for [Comunicación sensor]", err)
				SensorError = true // no notification sent
			}
		}
	} else {
		if !SensorError {
			SensorError = true
			log.Println("calderadevice:: caldera Error ", SensorError, " ", err)
			if err := notifier.SendMail("Error de sensor "+time.Now().Format("15:04:05"), "No hay comunicación con el sensor de temperatura"); err != nil {
				log.Println("Cannot send mail for [Error de sensor]", err)
				SensorError = false // no notification sent
			}
		}
	}
}

// Sensors get all connected sensor IDs as array
func Sensors() ([]string, error) {
	data, err := ioutil.ReadFile("/sys/bus/w1/devices/w1_bus_master1/w1_master_slaves")
	if err != nil {
		return nil, err
	}

	sensors := strings.Split(string(data), "\n")
	if len(sensors) > 0 {
		sensors = sensors[:len(sensors)-1]
	}

	return sensors, nil
}

// Temperature get the temperature of a given sensor
func Temperature(sensor string) (float64, error) {
	data, err := ioutil.ReadFile("/sys/bus/w1/devices/" + sensor + "/w1_slave")
	if err != nil {
		return 0.0, errReadSensor
	}

	raw := string(data)

	i := strings.LastIndex(raw, "t=")
	if i == -1 {
		return 0.0, errReadSensor
	}

	c, err := strconv.ParseFloat(raw[i+2:len(raw)-1], 64)
	if err != nil {
		return 0.0, errReadSensor
	}

	return c / 1000.0, nil
}

// GetTemp returns temperature of the first sensor found
func GetTemp() (float64, error) {
	sensors, err := Sensors()
	if err != nil {
		return 0.0, err
	}
	for _, sensor := range sensors {
		t, err := Temperature(sensor)
		if err == nil {
			return t, nil
		}
	}
	return 0.0, errReadSensor
}
