package jobs

import (
	"homekit/app/calderadevice"
	"homekit/app/models"
	"homekit/app/tempsensor"
	"time"

	"github.com/revel/modules/jobs/app/jobs"
	gorp "github.com/revel/modules/orm/gorp/app"
	"github.com/revel/revel"
)

// MeassurementRegister insert meassurements into the db
type MeassurementRegister struct {
}

// Run runnable method of MeassurementRegister
func (c MeassurementRegister) Run() {
	sensorTemp, err := tempsensor.GetTemp()
	tempsensor.CheckError(err)
	calderaTemp, err := calderadevice.GetTemp()
	calderadevice.CheckError(err)
	meassurement := models.Meassurement{
		Time:       time.Now().Unix(),
		ValCaldera: calderaTemp,
		ValSensor:  sensorTemp,
		Active:     calderadevice.CalderaActive}
	if err := gorp.Db.Map.Insert(&meassurement); err != nil {
		panic(err)
	}
}

func init() {
	revel.OnAppStart(func() {
		jobs.Now(MeassurementRegister{})
		jobs.Schedule("@every 10m", MeassurementRegister{})
	})
}
