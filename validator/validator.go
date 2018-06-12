package validator

import (
	"errors"
	"reflect"
	"strings"
)

func IsZeroOfUnderlyingType(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

func ValidateRequest(req interface{}) error {
	v := reflect.ValueOf(req)
	var missingFields []string
	for i := 0; i < v.Type().NumField(); i++ {
		json := v.Type().Field(i).Tag.Get("json")
		valid := v.Type().Field(i).Tag.Get("valid")
		if valid == "required" && IsZeroOfUnderlyingType(v.Field(i).Interface()) {
			missingFields = append(missingFields, json)
		}
	}
	if len(missingFields) > 0 {
		missing := strings.Join(missingFields, ",")
		return errors.New("Missing fields " + missing)
	}
	return nil
}
