package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/BabyBoChen/bbljfooddiary/models"
	"github.com/BabyBoChen/bbljfooddiary/services"
	"github.com/BabyBoChen/bbljfooddiary/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	_ "github.com/lib/pq"
)

var sessionStore *session.Store

func main() {
	envVars := services.ReadEnvironmentVariables()
	engine := html.New("./views", ".html")
	// Reload the templates on each render, good for development
	engine.Reload(envVars.Config == "debug") // Optional. Default: false
	// Debug will print each template that is parsed, good for debugging
	engine.Debug(envVars.Config == "debug") // Optional. Default: false

	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			vm := make(models.ErrorViewModel)
			vm["ErrorMessage"] = err
			return ctx.Render("/error", vm)
		},
	})

	sessionStore = session.New()
	AuthorizationMiddleware(app)
	app.Static("/", "./wwwroot")
	app.Get("/", index)
	app.Get("/newCuisine", newCuisine)
	app.Post("/newCuisine", postNewCuisine)
	app.Get("/cuisineList", cuisineList)
	app.Post("/query", query)
	app.Get("/login", login)
	app.Post("/login", postLogin)
	app.Get("/editCuisine", editCuisine)
	app.Post("/editCuisine", postEditCuisine)
	app.Get("/deleteCuisine", deleteCuisine)
	app.Get("/clearCache", clearCache)
	app.Get("/errorPage", errorPage)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", envVars.Port)))
}

func index(c *fiber.Ctx) error {
	viewModel := models.NewIndexViewModel()
	return c.Render("index", viewModel)
}

func newCuisine(c *fiber.Ctx) error {
	return c.Render("newCuisine", nil)
}

func postNewCuisine(c *fiber.Ctx) error {
	form, err := c.MultipartForm()

	var cuisineName string
	var unitPrice float64
	var cuisineType int64
	var isOneSet bool
	var lastOrderDate time.Time
	var review int64
	var restaurant string
	var address string
	var remark string
	if err == nil {
		cuisineName = utils.GetValueFromFormData(form.Value, "CuisineName")
		unitPrice_str := utils.GetValueFromFormData(form.Value, "UnitPrice")
		cuisineType_str := utils.GetValueFromFormData(form.Value, "CuisineType")
		isOneSet_str := utils.GetValueFromFormData(form.Value, "IsOneSet")
		lastOrderDate_str := utils.GetValueFromFormData(form.Value, "LastOrderDate")
		review_str := utils.GetValueFromFormData(form.Value, "Review")
		restaurant = utils.GetValueFromFormData(form.Value, "Restaurant")
		address = utils.GetValueFromFormData(form.Value, "Address")
		remark = utils.GetValueFromFormData(form.Value, "Remark")

		var errType error
		errMsg := ""
		unitPrice, errType = strconv.ParseFloat(unitPrice_str, 64)
		if errType != nil {
			errMsg += errType.Error() + ";"
		}
		cuisineType, errType = strconv.ParseInt(cuisineType_str, 10, 32)
		if errType != nil {
			errMsg += errType.Error() + ";"
		}
		if isOneSet_str == "1" {
			isOneSet = true
		} else {
			isOneSet = false
		}
		lastOrderDate, errType = time.Parse("2006-01-02", lastOrderDate_str)
		if errType != nil {
			errMsg += errType.Error() + ";"
		}
		review, errType = strconv.ParseInt(review_str, 10, 32)
		if errType != nil {
			errMsg += errType.Error() + ";"
		}
		if len(errMsg) > 0 {
			err = errors.New(errMsg)
		}
	}

	var header *multipart.FileHeader
	var cuisineImage io.Reader = nil
	if err == nil {
		var noImageError error
		header, noImageError = c.FormFile("CuisineImage")
		if noImageError == nil {
			cuisineImage, err = header.Open()
		}
	}

	var service *services.CuisineService
	if err == nil {
		service, err = services.NewCuisineService()
	}

	if err == nil {
		defer service.Dispose()
		newCuisine := make(map[string]interface{})
		newCuisine["cuisine_name"] = cuisineName
		newCuisine["unit_price"] = unitPrice
		newCuisine["last_order_date"] = lastOrderDate
		newCuisine["review"] = review
		newCuisine["restaurant"] = restaurant
		newCuisine["address"] = address
		newCuisine["remark"] = remark
		newCuisine["is_one_set"] = isOneSet
		newCuisine["cuisine_type"] = cuisineType
		err = service.SaveNewCuisine(newCuisine, cuisineImage)
	}

	if err == nil {
		return c.Redirect("/")
	} else {
		vm := make(models.ErrorViewModel)
		vm["ErrorMessage"] = err
		return c.Render("error", vm)
	}
}

func cuisineList(c *fiber.Ctx) error {
	vm := models.NewCuisineListViewModel()
	return c.Render("cuisineList", vm)
}

func query(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	var jsonBody string
	if err == nil {
		formData := make(map[string]string)
		formData["IsKeyword"] = utils.GetValueFromFormData(form.Value, "IsKeyword")
		formData["Keyword"] = utils.GetValueFromFormData(form.Value, "Keyword")
		formData["CuisineName"] = utils.GetValueFromFormData(form.Value, "CuisineName")
		formData["UnitPriceOrder"] = utils.GetValueFromFormData(form.Value, "UnitPriceOrder")
		formData["CuisineType"] = utils.GetValueFromFormData(form.Value, "CuisineType")
		formData["LastOrderDate"] = utils.GetValueFromFormData(form.Value, "LastOrderDate")
		formData["ReviewOrder"] = utils.GetValueFromFormData(form.Value, "ReviewOrder")
		formData["Restaurant"] = utils.GetValueFromFormData(form.Value, "Restaurant")
		formData["Address"] = utils.GetValueFromFormData(form.Value, "Address")
		formData["Remark"] = utils.GetValueFromFormData(form.Value, "Remark")
		vm := models.NewQueryViewModel(formData)
		jsonBody, err = vm.Result()
	}

	if err == nil {
		return c.SendString(jsonBody)
	} else {
		return c.SendString(err.Error())
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

func editCuisine(c *fiber.Ctx) error {
	var vm models.EditCuisineViewModel
	var err error

	q := c.Queries()
	if utils.MapContainsKey(q, "id") {
		vm, err = models.NewEditCuisineViewModel(q["id"])
	} else {
		err = errors.New("not found")
	}

	if err == nil {
		return c.Render("editCuisine", vm)
	} else {
		vm := make(models.ErrorViewModel)
		vm["ErrorMessage"] = err
		return c.Render("error", vm)
	}
}

func postEditCuisine(c *fiber.Ctx) error {
	form, err := c.MultipartForm()

	var cuisineId int64
	var cuisineName string
	var unitPrice float64
	var cuisineType int64
	var isOneSet bool
	var lastOrderDate time.Time
	var review int64
	var restaurant string
	var address string
	var remark string
	if err == nil {
		cuisineId_str := utils.GetValueFromFormData(form.Value, "CuisineId")
		cuisineName = utils.GetValueFromFormData(form.Value, "CuisineName")
		unitPrice_str := utils.GetValueFromFormData(form.Value, "UnitPrice")
		cuisineType_str := utils.GetValueFromFormData(form.Value, "CuisineType")
		isOneSet_str := utils.GetValueFromFormData(form.Value, "IsOneSet")
		lastOrderDate_str := utils.GetValueFromFormData(form.Value, "LastOrderDate")
		review_str := utils.GetValueFromFormData(form.Value, "Review")
		restaurant = utils.GetValueFromFormData(form.Value, "Restaurant")
		address = utils.GetValueFromFormData(form.Value, "Address")
		remark = utils.GetValueFromFormData(form.Value, "Remark")

		var errType error
		errMsg := ""
		cuisineId, errType = strconv.ParseInt(cuisineId_str, 10, 64)
		if errType != nil {
			errMsg += errType.Error() + ";"
		}
		unitPrice, errType = strconv.ParseFloat(unitPrice_str, 64)
		if errType != nil {
			errMsg += errType.Error() + ";"
		}
		cuisineType, errType = strconv.ParseInt(cuisineType_str, 10, 32)
		if errType != nil {
			errMsg += errType.Error() + ";"
		}
		if isOneSet_str == "1" {
			isOneSet = true
		} else {
			isOneSet = false
		}
		lastOrderDate, errType = time.Parse("2006-01-02", lastOrderDate_str)
		if errType != nil {
			errMsg += errType.Error() + ";"
		}
		review, errType = strconv.ParseInt(review_str, 10, 32)
		if errType != nil {
			errMsg += errType.Error() + ";"
		}
		if len(errMsg) > 0 {
			err = errors.New(errMsg)
		}
	}

	var header *multipart.FileHeader
	var cuisineImage io.Reader = nil
	if err == nil {
		var noImageError error
		header, noImageError = c.FormFile("CuisineImage")
		if noImageError == nil {
			cuisineImage, err = header.Open()
		}
	}

	var service *services.CuisineService
	if err == nil {
		service, err = services.NewCuisineService()
	}

	if err == nil {
		defer service.Dispose()
		newCuisine := make(map[string]interface{})
		newCuisine["cuisine_id"] = cuisineId
		newCuisine["cuisine_name"] = cuisineName
		newCuisine["unit_price"] = unitPrice
		newCuisine["last_order_date"] = lastOrderDate
		newCuisine["review"] = review
		newCuisine["restaurant"] = restaurant
		newCuisine["address"] = address
		newCuisine["remark"] = remark
		newCuisine["is_one_set"] = isOneSet
		newCuisine["cuisine_type"] = cuisineType
		err = service.SaveCuisine(newCuisine, cuisineImage)
	}

	if err == nil {
		return c.Redirect("/")
	} else {
		vm := make(models.ErrorViewModel)
		vm["ErrorMessage"] = err
		return c.Render("error", vm)
	}
}

func deleteCuisine(c *fiber.Ctx) error {
	var err error
	q := c.Queries()
	if utils.MapContainsKey(q, "id") {

	} else {
		err = errors.New("not found")
	}

	var service *services.CuisineService
	if err == nil {
		service, err = services.NewCuisineService()
	}

	if err == nil {
		defer service.Dispose()
		err = service.DeleteCuisine(q["id"])
	}

	if err == nil {
		return c.Redirect("/")
	} else {
		vm := make(models.ErrorViewModel)
		vm["ErrorMessage"] = err
		return c.Render("error", vm)
	}
}

func clearCache(c *fiber.Ctx) error {
	services.ClearCache()

	return c.Redirect("/")
}

func errorPage(c *fiber.Ctx) error {
	vm := make(models.ErrorViewModel)
	vm["ErrorMessage"] = "Error page..."
	return c.Render("error", vm)
}
