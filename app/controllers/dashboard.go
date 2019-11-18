package controllers

import (
	"homekit/app/routes"
	"log"

	"github.com/revel/revel"
)

type Dashboard struct {
	Application
}

func (c Dashboard) checkUser() revel.Result {
	log.Println("checkuser")
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Application.Index())
	}
	return nil
}

func (c Dashboard) Index() revel.Result {
	c.Log.Info("Fetching index")
	return c.Render()
}
