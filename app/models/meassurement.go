package models

// Meassurement of the system
type Meassurement struct {
	ValCaldera  float64
	ValSensor   float64
	ValExterior float64
	Active      bool
	Time        int64
}
