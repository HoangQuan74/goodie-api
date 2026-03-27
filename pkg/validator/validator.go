package validator

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	apperrors "github.com/HoangQuan74/goodie-api/pkg/errors"
)

func Init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("phone_vn", validateVNPhone)
	}
}

func validateVNPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if len(phone) < 10 || len(phone) > 11 {
		return false
	}
	return strings.HasPrefix(phone, "0") || strings.HasPrefix(phone, "+84")
}

func FormatValidationErrors(err error) *apperrors.AppError {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		messages := make([]string, 0, len(validationErrors))
		for _, e := range validationErrors {
			messages = append(messages, fmt.Sprintf("field '%s' failed on '%s' validation", e.Field(), e.Tag()))
		}
		return apperrors.BadRequest(strings.Join(messages, "; "))
	}
	return apperrors.BadRequest("invalid request body")
}
