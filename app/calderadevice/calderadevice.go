package calderadevice

import (
	"homekit/app/notifier"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	// Caldera caldera device name
	Caldera       string = "caldera"
	domain        string = ".local"
	uptimeURI     string = "/uptime"
	tempURI       string = "/temp"
	relayStateURI string = "/relay?state="
)

var (
	// CalderaError state of error of the device
	CalderaError bool = false
	// CalderaActive state of the caldera output
	CalderaActive bool = false
)

// CheckError check if error state and notify
func CheckError(err error) {
	if err == nil {
		if CalderaError {
			CalderaError = false
			if err := notifier.SendMail("Comunicación establecida "+time.Now().Format("15:04:05"), "Se ha establecido comunicación con el dispositivo caldera"); err != nil {
				log.Println("Cannot send mail for [Comunicación establecida]", err)
				CalderaError = true // no notification sent
			}
		}
	} else {
		if !CalderaError {
			CalderaError = true
			if err := notifier.SendMail("Error de comunicación "+time.Now().Format("15:04:05"), "No hay comunicación con el dispositivo caldera"); err != nil {
				log.Println("Cannot send mail for [Error de comunicación]", err)
				CalderaError = false // no notification sent
			}
		}
	}
}

// GetTemp returns the caldera temperature readed
func GetTemp() (float64, error) {
	resp, err := http.Get("http://" + Caldera + domain + tempURI)
	if err != nil {
		return 0.0, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	temp, err := strconv.ParseFloat(string(body), 64)
	if err != nil {
		return 0.0, err
	}
	return temp, nil
}

// GetUpTime returns the caldera up time
func GetUpTime() (string, error) {
	resp, err := http.Get("http://" + Caldera + domain + uptimeURI)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

// SetState sets the caldera state
func SetState(state bool) error {
	active := "off"
	if state {
		active = "on"
	}
	resp, err := http.Get("http://" + Caldera + domain + relayStateURI + active)
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
