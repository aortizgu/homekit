// GENERATED CODE - DO NOT EDIT
// This file provides a way of creating URL's based on all the actions
// found in all the controllers.
package routes

import "github.com/revel/revel"


type tJobs struct {}
var Jobs tJobs


func (_ tJobs) Status(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Jobs.Status", args).URL
}


type tController struct {}
var Controller tController


func (_ tController) Begin(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Controller.Begin", args).URL
}

func (_ tController) Rollback(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Controller.Rollback", args).URL
}

func (_ tController) Commit(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Controller.Commit", args).URL
}


type tStatic struct {}
var Static tStatic


func (_ tStatic) Serve(
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.Serve", args).URL
}

func (_ tStatic) ServeDir(
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.ServeDir", args).URL
}

func (_ tStatic) ServeModule(
		moduleName string,
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "moduleName", moduleName)
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.ServeModule", args).URL
}

func (_ tStatic) ServeModuleDir(
		moduleName string,
		prefix string,
		filepath string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "moduleName", moduleName)
	revel.Unbind(args, "prefix", prefix)
	revel.Unbind(args, "filepath", filepath)
	return revel.MainRouter.Reverse("Static.ServeModuleDir", args).URL
}


type tTestRunner struct {}
var TestRunner tTestRunner


func (_ tTestRunner) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.Index", args).URL
}

func (_ tTestRunner) Suite(
		suite string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "suite", suite)
	return revel.MainRouter.Reverse("TestRunner.Suite", args).URL
}

func (_ tTestRunner) Run(
		suite string,
		test string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "suite", suite)
	revel.Unbind(args, "test", test)
	return revel.MainRouter.Reverse("TestRunner.Run", args).URL
}

func (_ tTestRunner) List(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("TestRunner.List", args).URL
}


type tApplication struct {}
var Application tApplication


func (_ tApplication) AddUser(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Application.AddUser", args).URL
}

func (_ tApplication) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Application.Index", args).URL
}

func (_ tApplication) Login(
		username string,
		password string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "username", username)
	revel.Unbind(args, "password", password)
	return revel.MainRouter.Reverse("Application.Login", args).URL
}

func (_ tApplication) Logout(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Application.Logout", args).URL
}


type tDashboard struct {}
var Dashboard tDashboard


func (_ tDashboard) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Dashboard.Index", args).URL
}

func (_ tDashboard) Live(
		user string,
		ws interface{},
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "user", user)
	revel.Unbind(args, "ws", ws)
	return revel.MainRouter.Reverse("Dashboard.Live", args).URL
}


type tRules struct {}
var Rules tRules


func (_ tRules) Index(
		) string {
	args := make(map[string]string)
	
	return revel.MainRouter.Reverse("Rules.Index", args).URL
}

func (_ tRules) GetDeviceRules(
		device string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "device", device)
	return revel.MainRouter.Reverse("Rules.GetDeviceRules", args).URL
}

func (_ tRules) SetDeviceRules(
		device string,
		) string {
	args := make(map[string]string)
	
	revel.Unbind(args, "device", device)
	return revel.MainRouter.Reverse("Rules.SetDeviceRules", args).URL
}


