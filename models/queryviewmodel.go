package models

import (
	"encoding/json"
	"strconv"

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
	var totalCount int = 0

	if vm.formData["IsKeyword"] == "1" {
		service, err = services.NewCuisineService()
		if err == nil {
			defer service.Dispose()
			result, totalCount, err = service.QueryWithKeyword(vm.formData)
		}
	} else {
		service, err = services.NewCuisineService()
		if err == nil {
			defer service.Dispose()
			result, totalCount, err = service.QueryWithFields(vm.formData)
		}
	}

	pageSize := 20
	pageSize, _ = strconv.Atoi(vm.formData["size"])

	var jb []byte
	if err == nil {
		jsonObj := make(map[string]interface{})
		jsonObj["last_page"] = calculateLastPage(totalCount, pageSize)
		jsonObj["data"] = result
		jb, err = json.Marshal(jsonObj)
	}

	if err == nil {
		jsonBody = string(jb)
	}

	return jsonBody, err
}

func calculateLastPage(totalCount, pageSize int) int {
	if pageSize <= 0 {
		return 0 // Or handle as an error, depending on your desired behavior for invalid pageSize
	}

	lastPage := totalCount / pageSize
	if totalCount%pageSize != 0 {
		lastPage++
	}
	return lastPage
}
