package models

import (
	"encoding/json"

	"github.com/BabyBoChen/bbljfooddiary/services"
)

type QueryViewModel struct {
	formData map[string]string
}

func NewQueryViewModel(formData map[string]string) *QueryViewModel {
	var vm QueryViewModel
	vm.formData = formData
	return &vm
}

func (vm *QueryViewModel) Result() (string, error) {
	jsonBody := "{\"allCuisines:\"[],\"formData\":\"\"}"
	var err error

	var service *services.CuisineService
	var result []map[string]interface{}
	if vm.formData["IsKeyword"] == "1" {
		service, err = services.NewCuisineService()
		if err == nil {
			defer service.Dispose()
			result, err = service.QueryWithKeyword(vm.formData["Keyword"])
		}
	} else {
		service, err = services.NewCuisineService()
		if err == nil {
			defer service.Dispose()
			result, err = service.QueryWithFields(vm.formData)
		}
	}

	var jb []byte
	if err == nil {
		jsonObj := make(map[string]interface{})
		jsonObj["allCuisines"] = result
		jsonObj["formData"] = vm.formData
		jb, err = json.Marshal(jsonObj)
	}

	if err == nil {
		jsonBody = string(jb)
	}

	return jsonBody, err
}
