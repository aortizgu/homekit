package jobs

import (
	"github.com/revel/modules/jobs/app/jobs"
	gorp "github.com/revel/modules/orm/gorp/app"
	"github.com/revel/revel"
)

const (
	maxMeassurements int64 = (24 * 60) / 10
)

// DbLimiter limits the entries for the db
type DbLimiter struct {
}

// Run runnable method of DbLimiter
func (c DbLimiter) Run() {
	count, err := gorp.Db.Map.SelectInt("select count(*) from Meassurement")
	if err == nil {
		toDelete := count - maxMeassurements
		if toDelete > 0 {
			gorp.Db.Map.Exec("delete from Meassurement where Time IN (SELECT Time from Meassurement order by Time asc limit ?)", toDelete)
		}
	}
}

func init() {
	revel.OnAppStart(func() {
		jobs.Now(DbLimiter{})
		jobs.Schedule("@every 10s", DbLimiter{})
	})
}
