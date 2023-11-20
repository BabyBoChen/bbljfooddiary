package models

import (
	"encoding/json"

	"github.com/BabyBoChen/bbljfooddiary/services"
)

type IndexViewModel = map[string]interface{}

func NewIndexViewModel() IndexViewModel {
	vm := make(IndexViewModel)
	vm["Top10Main"] = "[]"
	vm["Top10Dessert"] = "[]"
	vm["Top10Buffet"] = "[]"
	vm["Top10NewCuisine"] = "[]"

	ds := make(services.DataSet)
	service, err := services.NewCuisineService()
	if err == nil {
		defer service.Dispose()
		ds, err = service.GetTop10Cuisines()
	}
	if err == nil {
		var jsonArr []byte
		jsonArr, _ = json.Marshal(ds["Top10Main"])
		vm["Top10Main"] = string(jsonArr)

		jsonArr, _ = json.Marshal(ds["Top10Dessert"])
		vm["Top10Dessert"] = string(jsonArr)

		jsonArr, _ = json.Marshal(ds["Top10Buffet"])
		vm["Top10Buffet"] = string(jsonArr)

		jsonArr, _ = json.Marshal(ds["Top10NewCuisine"])
		vm["Top10NewCuisine"] = string(jsonArr)
	}
	return vm
}
