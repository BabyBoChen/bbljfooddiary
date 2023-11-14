package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/BabyBoChen/bbljfooddiary/models"
	"github.com/BabyBoChen/bbljfooddiary/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	_ "github.com/lib/pq"
)

var sessionStore *session.Store

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
		cuisineName = form.Value["CuisineName"][0]
		unitPrice_str := form.Value["UnitPrice"][0]
		cuisineType_str := form.Value["CuisineType"][0]
		isOneSet_str := form.Value["IsOneSet"][0]
		lastOrderDate_str := form.Value["LastOrderDate"][0]
		review_str := form.Value["Review"][0]
		restaurant = form.Value["Restaurant"][0]
		address = form.Value["Address"][0]
		remark = form.Value["Remark"][0]

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

	var service services.CuisineService
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
		err = service.SaveNewCuisine(newCuisine)

		// fmt.Println(cuisineName)
		// fmt.Println(unitPrice)
		// fmt.Println(cuisineType)
		// fmt.Println(isOneSet)
		// fmt.Println(lastOrderDate)
		// fmt.Println(review)
		// fmt.Println(restaurant)
		// fmt.Println(address)
		// fmt.Println(remark)
	}

	if err == nil {
		return c.Redirect("/")
	} else {
		vm := make(models.ErrorViewModel)
		vm["ErrorMessage"] = err
		return c.Render("/error", vm)
	}
}

func cuisineList(c *fiber.Ctx) error {
	vm := models.NewCuisineListViewModel()
	return c.Render("cuisineList", vm)
}

func main() {
	envVars := services.ReadEnvironmentVariables()
	engine := html.New("./views", ".html")
	// Reload the templates on each render, good for development
	engine.Reload(envVars.Config == "debug") // Optional. Default: false
	// Debug will print each template that is parsed, good for debugging
	engine.Debug(envVars.Config == "debug") // Optional. Default: false

	sessionStore = session.New()

	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			vm := make(models.ErrorViewModel)
			vm["ErrorMessage"] = err
			return ctx.Render("/error", vm)
		},
	})

	app.Static("/", "./wwwroot")

	app.Get("/", index)
	app.Get("/newCuisine", newCuisine)
	app.Post("/newCuisine", postNewCuisine)
	app.Get("/cuisineList", cuisineList)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", envVars.Port)))
}
