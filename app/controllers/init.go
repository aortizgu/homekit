package controllers

import "github.com/revel/revel"

func init() {
	revel.InterceptMethod(Application.AddUser, revel.BEFORE)
	revel.InterceptMethod(Dashboard.checkUser, revel.BEFORE)
	revel.InterceptMethod(Rules.checkUser, revel.BEFORE)
}
