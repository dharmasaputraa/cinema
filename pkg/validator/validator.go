package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Validator struct {
	v *validator.Validate
}

func New() *Validator {
	v := validator.New()

	// pakai json tag sebagai field name
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// custom rules
	v.RegisterValidation("uuid", isUUID)
	v.RegisterValidation("phone_id", isIndonesianPhone)

	return &Validator{v}
}

func (val *Validator) Validate(s interface{}) map[string][]string {
	err := val.v.Struct(s)
	if err == nil {
		return nil
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string][]string{
			"_": {err.Error()},
		}
	}

	fields := make(map[string][]string)

	for _, e := range errs {
		field := buildJSONPath(e)
		fields[field] = append(fields[field], humanize(e))
	}

	return fields
}

// ── Helpers ─────────────────────────────────────────────────

func buildJSONPath(e validator.FieldError) string {
	namespace := e.StructNamespace() // User.Address.Street

	parts := strings.Split(namespace, ".")
	if len(parts) <= 1 {
		return e.Field()
	}

	// buang root struct
	parts = parts[1:]

	return strings.Join(parts, ".")
}

// ── Custom Rules ────────────────────────────────────────────

func isUUID(fl validator.FieldLevel) bool {
	_, err := uuid.Parse(fl.Field().String())
	return err == nil
}

func isIndonesianPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	if strings.HasPrefix(phone, "+62") {
		phone = strings.TrimPrefix(phone, "+62")
	} else if strings.HasPrefix(phone, "08") {
		phone = strings.TrimPrefix(phone, "08")
	} else {
		return false
	}

	if len(phone) < 8 || len(phone) > 13 {
		return false
	}

	for _, c := range phone {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

// ── Error Messages ──────────────────────────────────────────

func humanize(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", e.Field())
	case "phone_id":
		return fmt.Sprintf("%s must be a valid Indonesian phone number", e.Field())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}
