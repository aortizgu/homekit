package led

import (
	"fmt"
	"homekit/app/msgbroker"
	"os"
)

const (
	line1 = 0
	line2 = 1
)

var (
	files = []string{"/tmp/led_1.txt", "/tmp/led_2.txt"}
)

func setLine(i int, msg string) {
	f, err := os.Create(files[i])
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(msg)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Update(e msgbroker.Event) {
	setLine(line1, "In "+fmt.Sprintf("%.1f", e.SensorTemp)+" Ex "+fmt.Sprintf("%.1f", e.ExternalTemp))
	msg := ""
	if !e.Status {
		msg = "Error en sistema"
	} else {
		if e.Active && e.Manual {
			msg = "ON Manual"
		} else if e.Active && !e.Manual {
			msg = "ON Automatico"
		} else {
			msg = "OFF"
		}
	}
	setLine(line2, msg)
}
