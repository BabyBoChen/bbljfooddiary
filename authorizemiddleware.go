package main

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AuthorizationMiddleware(app *fiber.App) {
	app.Use([]string{"/newCuisine", "/saveCuisine", "/deleteCuisine", "/clearCache"}, authorize)
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
		sess.Set("redirectUrl", c.Request().URI().String())
		sess.SetExpiry(time.Minute * 15)
		sess.Save()
		return c.Redirect("/login")
	}
}

func IsLogin(c *fiber.Ctx) bool {
	isAuthorized := false
	sess, err := sessionStore.Get(c)
	if err == nil {
		isAuthorizedSess := sess.Get("isAuthorized")
		if isAuthorizedSess == true {
			isAuthorized = true
		}
	}
	return isAuthorized
}
