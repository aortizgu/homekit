package jobs

import (
	"fmt"
	"homekit/app/models"

	"github.com/revel/revel"

	"github.com/revel/modules/jobs/app/jobs"
	gorp "github.com/revel/modules/orm/gorp/app"
)

// Periodically count the bookings in the database.
type UserCounter struct{}

func (c UserCounter) Run() {
	users, err := gorp.Db.Map.Select(models.User{},
		`select * from User`)
	if err != nil {
		panic(err)
	}
	fmt.Printf("There are %d users.\n", len(users))
}

func init() {
	revel.OnAppStart(func() {
		jobs.Schedule("@every 10s", UserCounter{})
	})
}
