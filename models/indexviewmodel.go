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
	var top10Main []map[string]interface{}
	var top10Dessert []map[string]interface{}
	var top10Buffet []map[string]interface{}

	service, err := services.NewCuisineService()
	if err == nil {
		defer service.Dispose()
		top10Main, top10Dessert, top10Buffet, err = service.GetTop10Cuisines()
	}
	if err == nil {
		var jsonArr []byte
		jsonArr, _ = json.Marshal(top10Main)
		vm["Top10Main"] = string(jsonArr)

		jsonArr, _ = json.Marshal(top10Dessert)
		vm["Top10Dessert"] = string(jsonArr)

		jsonArr, _ = json.Marshal(top10Buffet)
		vm["Top10Buffet"] = string(jsonArr)
	}
	return vm
}
