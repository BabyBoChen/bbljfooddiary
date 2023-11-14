package models

import (
	"encoding/json"

	"github.com/BabyBoChen/bbljfooddiary/services"
)

type CuisineListViewModel = map[string]interface{}

func NewCuisineListViewModel() CuisineListViewModel {
	vm := make(CuisineListViewModel)
	vm["AllCuisine"] = "[]"
	var allCuisine []map[string]interface{}

	service, err := services.NewCuisineService()

	if err == nil {
		defer service.Dispose()
		allCuisine, err = service.ListAllCuisine()
	}

	if err == nil {
		var jsonArr []byte
		jsonArr, _ = json.Marshal(allCuisine)
		vm["AllCuisine"] = string(jsonArr)
	}
	return vm
}
