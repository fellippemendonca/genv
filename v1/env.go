package v1

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// Config example struct
type Example struct {
	Port string `name:"SERVICE_NAME_PORT" required:"true"`
}

const NAME = "name"
const REQUIRED = "required"
const DEFAULTS = "default"

func Load(config interface{}) error {
	structFields := reflect.TypeOf(config).Elem()
	structValues := reflect.ValueOf(config).Elem()

	for i := 0; i < structFields.NumField(); i++ {
		field := structFields.Field(i)
		value := structValues.Field(i)

		envVariable := field.Tag.Get(NAME)
		envValue := os.Getenv(envVariable)

		required := field.Tag.Get(REQUIRED)
		defaults := field.Tag.Get(DEFAULTS)

		if envValue == "" && defaults != "" {
			envValue = field.Tag.Get(DEFAULTS)
		}

		if required == "true" && envValue == "" {
			return fmt.Errorf("required env var not set: %s", envVariable)
		}

		if err := setValue(value, envValue); err != nil {
			return fmt.Errorf("failed setting env var, %w", err)
		}
	}

	return nil
}

func setValue(value reflect.Value, envValue string) error {
	switch value.Kind() {
	case reflect.String:
		value.SetString(envValue)
	case reflect.Bool:
		v, err := strconv.ParseBool(envValue)
		if err != nil {
			return err
		}
		value.SetBool(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(envValue, 10, 64)
		if err != nil {
			return err
		}
		value.SetInt(v)
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(envValue, 64)
		if err != nil {
			return err
		}
		value.SetFloat(v)
	default:
		return fmt.Errorf("unsupported type: %s", value.Kind())
	}
	return nil
}
