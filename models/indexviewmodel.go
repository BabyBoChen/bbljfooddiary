package models

import (
	"encoding/json"

	"github.com/BabyBoChen/bbljfooddiary/services"
)

type IndexViewModel = map[string]interface{}

func NewIndexViewModel() IndexViewModel {
	vm := make(IndexViewModel)
	vm["Top10Cuisines"] = "[]"
	service, err := services.NewCuisineService()

	var top10 []map[string]interface{}
	if err == nil {
		defer service.Dispose()
		top10, err = service.GetTop10Cuisines()
	}
	if err == nil {
		json, _ := json.Marshal(top10)
		vm["Top10Cuisines"] = string(json)
	}
	return vm
}
