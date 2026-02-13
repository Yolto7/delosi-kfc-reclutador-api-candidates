package validators

import (
	"fmt"
	"reflect"
	"strings"

	errorCustom "github.com/Yolto7/api-candidates/pkg/domain/error"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
    name := fld.Tag.Get("json")
    if name == "-" {
      return ""
    }
    if idx := strings.Index(name, ","); idx != -1 {
      name = name[:idx]
    }

    return name
	})

	// Custom rule: notblank = no espacios en blanco
	validate.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		return strings.TrimSpace(fl.Field().String()) != ""
	})
}

func RegisterCustomValidator(tag string, fn validator.Func) error {
	return validate.RegisterValidation(tag, fn)
}

type FieldErrorDetail struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

func ValidateSchema(input any) error {
	trimStringFields(input)

	err := validate.Struct(input)
	if err == nil {
		return nil
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		issues := make([]FieldErrorDetail, 0, len(validationErrors))

		for _, ve := range validationErrors {
			issues = append(issues, FieldErrorDetail{
				Key:     ve.Field(),
				Message: fmt.Sprintf("validation failed on '%s' constraint", ve.Tag()),
			})
		}

		return errorCustom.New(
			errorCustom.BAD_REQUEST,
			"Invalid payload supplied",
			"ERR_INVALID_PAYLOAD",
			issues,
		)
	}

	return err
}

func trimStringFields(i any) {
  v := reflect.ValueOf(i)
  if v.Kind() == reflect.Ptr {
    v = v.Elem()
  }
  if v.Kind() != reflect.Struct {
    return
  }

  for i := 0; i < v.NumField(); i++ {
    f := v.Field(i)
    if !f.CanSet() {
      continue
    }

    switch f.Kind() {
    case reflect.String:
      f.SetString(strings.TrimSpace(f.String()))
    case reflect.Struct:
      trimStringFields(f.Addr().Interface())
    case reflect.Ptr:
      if f.Elem().Kind() == reflect.Struct {
        trimStringFields(f.Interface())
      }
    }
  }
}
