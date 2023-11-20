package models

import (
	"encoding/json"
	"time"

	"github.com/BabyBoChen/bbljfooddiary/services"
)

type EditCuisineViewModel = map[string]interface{}

func NewEditCuisineViewModel(cuisine_id string) (EditCuisineViewModel, error) {
	vm := make(EditCuisineViewModel)
	service, err := services.NewCuisineService()

	cuisine := make(map[string]interface{})
	if err == nil {
		defer service.Dispose()
		cuisine, err = service.GetCuisineByCuisineId(cuisine_id)
	}

	var cuisine_json []byte
	if err == nil {
		cuisine["cuisine_type1"] = ""
		cuisine["cuisine_type2"] = ""
		cuisine["cuisine_type3"] = ""
		if cuisine["cuisine_type"] == int64(1) {
			cuisine["cuisine_type1"] = "selected"
		} else if cuisine["cuisine_type"] == int64(2) {
			cuisine["cuisine_type2"] = "selected"
		} else if cuisine["cuisine_type"] == int64(3) {
			cuisine["cuisine_type3"] = "selected"
		}

		cuisine["is_one_set1"] = ""
		cuisine["is_one_set0"] = ""
		if cuisine["is_one_set"] == true {
			cuisine["is_one_set1"] = "checked"
		} else if cuisine["is_one_set"] == false {
			cuisine["is_one_set0"] = "checked"
		}

		lastOrderDate := cuisine["last_order_date"].(time.Time)
		cuisine["last_order_date"] = lastOrderDate.Format("2006-01-02")

		cuisine_json, err = json.Marshal(cuisine)
	}

	if err == nil {
		vm = cuisine
		vm["CuisineMap"] = string(cuisine_json)
	}

	return vm, err
}
