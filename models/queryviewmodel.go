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
	jsonBody := "[]"
	var err error

	var service *services.CuisineService
	var result []map[string]interface{}
	if vm.formData["IsKeyword"] == "1" {
		service, err = services.NewCuisineService()
		if err == nil {
			defer service.Dispose()
			result, err = service.QueryWithKeyword(vm.formData["Keyword"])
		}
	}

	var jb []byte
	if err == nil {
		jb, err = json.Marshal(result)
	}

	if err == nil {
		jsonBody = string(jb)
	}

	return jsonBody, err
}
