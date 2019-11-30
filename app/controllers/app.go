package controllers

import (
	"homekit/app/models"

	"golang.org/x/crypto/bcrypt"

	"homekit/app/routes"

	"database/sql"

	gorpController "github.com/revel/modules/orm/gorp/app/controllers"
	"github.com/revel/revel"
)

// Application main application to login and clients
type Application struct {
	gorpController.Controller
}

// AddUser adds user name to vier args
func (c Application) AddUser() revel.Result {
	if user := c.connected(); user != nil {
		c.ViewArgs["user"] = user
	}
	return nil
}

func (c Application) connected() *models.User {
	if c.ViewArgs["user"] != nil {
		return c.ViewArgs["user"].(*models.User)
	}
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username.(string))
	}

	return nil
}

func (c Application) getUser(username string) (user *models.User) {
	user = &models.User{}
	_, err := c.Session.GetInto("fulluser", user, false)
	if user.Username == username {
		return user
	}

	err = c.Txn.SelectOne(user, c.Db.SqlStatementBuilder.Select("*").From("User").Where("Username=?", username))
	if err != nil {
		if err != sql.ErrNoRows {
			count, _ := c.Txn.SelectInt(c.Db.SqlStatementBuilder.Select("count(*)").From("User"))
			c.Log.Error("Failed to find user", "user", username, "error", err, "count", count)
		}
		return nil
	}
	c.Session["fulluser"] = user
	return
}

// Index show index page
func (c Application) Index() revel.Result {
	return c.Render()
}

// Login login client
func (c Application) Login(username, password string) revel.Result {
	user := c.getUser(username)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
		if err == nil {
			c.Session["user"] = username
			c.Session.SetNoExpiration()
			return c.Redirect(routes.Dashboard.Index())
		}
	}

	c.Flash.Out["username"] = username
	c.Flash.Error("Login failed")
	return c.Redirect(routes.Application.Index())
}

// Logout logout client
func (c Application) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.Application.Index())
}
