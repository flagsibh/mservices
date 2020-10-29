package data

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator"
)

// ErrInvalidProductPath is an error message when the product path is not valid
var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /products/[id]")

// ErrGenericError is a generic error message returned by a server
type ErrGenericError struct {
	Message string `json:"message"`
}

// ErrGenericErrors is a collection of generic or validation errors.
type ErrGenericErrors struct {
	Messages []string `json:"messages"`
}

// Validation contains the validator to apply validation rules
type Validation struct {
	validate *validator.Validate
}

// ValidationError wraps the validators FieldError so we do not
// expose this to out code
type ValidationError struct {
	validator.FieldError
}

func (v ValidationError) Error() string {
	return fmt.Sprintf(
		"Key: '%s' Error: Field validation for '%s' failed on the '%s' tag",
		v.Namespace(),
		v.Field(),
		v.Tag(),
	)
}

// ValidationErrors is a collection of ValidationError
type ValidationErrors []ValidationError

// Errors converts the slice into a string slice
func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

// NewValidation creates a new validatioin definition/tags on a product
func NewValidation() *Validation {
	validate := validator.New()
	validate.RegisterValidation("sku", skuValidator)
	return &Validation{validate}
}

// Validate runs validation on the item
func (v *Validation) Validate(i interface{}) ValidationErrors {
	errs := v.validate.Struct(i).(validator.ValidationErrors)

	if len(errs) == 0 {
		return nil
	}

	var returnErrs []ValidationError
	for _, err := range errs {
		// cast the FieldError into our ValidationError and append to the slice
		ve := ValidationError{err.(validator.FieldError)}
		returnErrs = append(returnErrs, ve)
	}

	return returnErrs
}

func skuValidator(fl validator.FieldLevel) bool {
	// SKU is abs-xsc-jhd
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return (len(matches) == 1)
}
