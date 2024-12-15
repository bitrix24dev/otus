package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, err := range v {
		sb.WriteString(fmt.Sprintf("Field: %s, Error: %s\n", err.Field, err.Err))
	}
	return sb.String()
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return errors.New("input is not a struct")
	}

	var validationErrors ValidationErrors

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)
		tag := field.Tag.Get("validate")

		if tag == "" {
			continue
		}

		validators := strings.Split(tag, "|")
		for _, validator := range validators {
			if err := validateField(field.Name, value, validator); err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateField(fieldName string, value reflect.Value, validator string) error {
	parts := strings.SplitN(validator, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid validator format: %s", validator)
	}

	switch value.Kind() {
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			if err := validateSingleValue(fieldName, value.Index(i), parts[0], parts[1]); err != nil {
				return err
			}
		}
	default:
		return validateSingleValue(fieldName, value, parts[0], parts[1])
	}
	return nil
}

func validateSingleValue(fieldName string, value reflect.Value, validatorType, param string) error {
	switch validatorType {
	case "len":
		return validateLength(fieldName, value, param)
	case "regexp":
		return validateRegexp(fieldName, value, param)
	case "in":
		return validateIn(fieldName, value, param)
	case "min":
		return validateMin(fieldName, value, param)
	case "max":
		return validateMax(fieldName, value, param)
	default:
		return fmt.Errorf("unknown validator: %s", validatorType)
	}
}

func validateLength(fieldName string, value reflect.Value, param string) error {
	expectedLen, err := strconv.Atoi(param)
	if err != nil {
		return fmt.Errorf("invalid length parameter: %s", param)
	}

	if value.Kind() == reflect.String && len(value.String()) != expectedLen {
		return fmt.Errorf("length of %s must be %d", fieldName, expectedLen)
	}

	return nil
}

func validateRegexp(fieldName string, value reflect.Value, param string) error {
	re, err := regexp.Compile(param)
	if err != nil {
		return fmt.Errorf("invalid regexp: %s", param)
	}

	if !re.MatchString(value.String()) {
		return fmt.Errorf("%s does not match regexp %s", fieldName, param)
	}
	return nil
}

func validateIn(fieldName string, value reflect.Value, param string) error {
	options := strings.Split(param, ",")
	switch value.Kind() {
	case reflect.String:
		for _, option := range options {
			if value.String() == option {
				return nil
			}
		}
	case reflect.Int:
		val := int(value.Int())
		for _, option := range options {
			if opt, err := strconv.Atoi(option); err == nil && val == opt {
				return nil
			}
		}
	default:
		return fmt.Errorf("unsupported type for in validation: %s", value.Kind())
	}
	return fmt.Errorf("%s is not in %v", fieldName, options)
}

func validateMin(fieldName string, value reflect.Value, param string) error {
	minVal, err := strconv.Atoi(param)
	if err != nil {
		return fmt.Errorf("invalid min parameter: %s", param)
	}

	if value.Kind() == reflect.Int && int(value.Int()) < minVal {
		return fmt.Errorf("%s must be at least %d", fieldName, minVal)
	}
	return nil
}

func validateMax(fieldName string, value reflect.Value, param string) error {
	maxVal, err := strconv.Atoi(param)
	if err != nil {
		return fmt.Errorf("invalid max parameter: %s", param)
	}

	if value.Kind() == reflect.Int && int(value.Int()) > maxVal {
		return fmt.Errorf("%s must be at most %d", fieldName, maxVal)
	}
	return nil
}
