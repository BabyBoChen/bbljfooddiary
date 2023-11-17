package main

import (
	"errors"
	"time"

	"github.com/BabyBoChen/bbljfooddiary/models"
	"github.com/BabyBoChen/bbljfooddiary/services"
	"github.com/BabyBoChen/bbljfooddiary/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func AuthorizationMiddleware(app *fiber.App) {
	app.Use([]string{"/newCuisine"}, authorize)
	app.Get("/login", login)
	app.Post("/login", postLogin)
}

func authorize(c *fiber.Ctx) error {
	isAuthorized := false

	sess, err := sessionStore.Get(c)
	if err == nil {
		isAuthorizedSess := sess.Get("isAuthorized")
		if isAuthorizedSess == true {
			isAuthorized = true
		}
	}

	if !isAuthorized {
		err = errors.New("unauthorized")
	}

	if err == nil {
		sess.SetExpiry(time.Minute * 15)
		sess.Save()
		return c.Next()
	} else {
		sess.Set("redirectUrl", c.Path())
		sess.SetExpiry(time.Minute * 15)
		sess.Save()
		return c.Redirect("/login")
	}
}

func login(c *fiber.Ctx) error {
	return c.Render("login", nil)
}

func postLogin(c *fiber.Ctx) error {
	form, err := c.MultipartForm()

	password := ""
	isAuthorized := false
	var sess *session.Session
	if err == nil {
		password = utils.GetValueFromFormData(form.Value, "Password")
		envVars := services.ReadEnvironmentVariables()
		if password == envVars.Password {
			isAuthorized = true
		}
		sess, err = sessionStore.Get(c)
	}

	redirectUrl := "/"
	if err == nil {
		sess.Set("isAuthorized", isAuthorized)
		if isAuthorized {
			if sess.Get("redirectUrl") != nil {
				redirectUrl = sess.Get("redirectUrl").(string)
				sess.Set("redirectUrl", nil)
			}
			sess.SetExpiry(time.Minute * 15)
			err = sess.Save()
		} else {
			err = errors.New("unauthorized")
		}
	}

	if err == nil {
		return c.Redirect(redirectUrl)
	} else {
		vm := make(models.ErrorViewModel)
		vm["ErrorMessage"] = err
		return c.Render("error", vm)
	}
}
