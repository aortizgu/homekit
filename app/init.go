package app

import (
	"homekit/app/models"
	"homekit/app/notifier"
	"log"
	"time"

	rgorp "github.com/revel/modules/orm/gorp/app"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gorp.v2"
)

func checkUserExists(db *gorp.DbMap, name, user, password string) {
	var u models.User
	err := db.SelectOne(&u, "select * from user where username=?", user)
	if err != nil {
		bcryptPassword, _ := bcrypt.GenerateFromPassword(
			[]byte(password), bcrypt.DefaultCost)
		user := &models.User{0, name, user, password, bcryptPassword}
		if err := db.Insert(user); err != nil {
			panic(err)
		}
	}

}

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	revel.OnAppStart(func() {
		Dbm := rgorp.Db.Map
		setColumnSizes := func(t *gorp.TableMap, colSizes map[string]int) {
			for col, size := range colSizes {
				t.ColMap(col).MaxSize = size
			}
		}

		t := Dbm.AddTable(models.User{}).SetKeys(true, "UserId")
		t.ColMap("Password").Transient = true
		setColumnSizes(t, map[string]int{
			"Username": 20,
			"Name":     100,
		})

		t = Dbm.AddTable(models.Meassurement{})

		//rgorp.Db.TraceOn(revel.AppLog)
		Dbm.CreateTables()
		checkUserExists(Dbm, "User demo", "demo", "demo")
		checkUserExists(Dbm, "Adrián Ortiz Gutiérrez", "aortiz", "orgut")
		if err := notifier.SendMail("Aplicación Arrancada "+time.Now().Format("15:04:05"), "Se acaba de iniciar la aplicación"); err != nil {
			log.Println("Cannot send mail for [Aplicación Arrancada]", err)
		}
	}, 5)
}

var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
