package controllers

import (
	"homekit/app/routes"
	"log"

	"github.com/revel/revel"
)

type Rules struct {
	Application
}

func (c Rules) checkUser() revel.Result {
	log.Println("checkuser")
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

func (c Rules) Index() revel.Result {
	c.Log.Info("Fetching index")
	return c.Render()
}
