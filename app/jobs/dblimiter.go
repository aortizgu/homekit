package jobs

import (
	"log"

	"github.com/revel/modules/jobs/app/jobs"
	gorp "github.com/revel/modules/orm/gorp/app"
	"github.com/revel/revel"
)

const (
	maxMeassurements int64 = 10 * 1000
)

type DbLimiter struct {
}

func (c DbLimiter) Run() {
	count, err := gorp.Db.Map.SelectInt("select count(*) from Meassurement")
	if err == nil {
		toDelete := count - maxMeassurements
		if toDelete > 0 {
			gorp.Db.Map.Exec("delete from Meassurement where Time IN (SELECT Time from Meassurement order by Time asc limit ?)", toDelete)
			log.Println("deleted oldest", toDelete, "entries")
		} else {
			log.Println(count, "entries in Meassurement table")
		}
	}
}

func init() {
	revel.OnAppStart(func() {
		jobs.Schedule("@every 10s", DbLimiter{})
	})
}
