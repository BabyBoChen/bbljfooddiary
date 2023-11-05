package utils

import (
	"errors"
	"reflect"
)

func InterfaceSlice(slice interface{}) ([]interface{}, error) {
	var arr []interface{}
	var err error
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		err = errors.New("InterfaceSlice() given a non-slice type")
	}

	if err == nil {
		if !s.IsNil() {
			arr = make([]interface{}, s.Len())
			for i := 0; i < s.Len(); i++ {
				arr[i] = s.Index(i).Interface()
			}
		}
	}

	return arr, err
}
