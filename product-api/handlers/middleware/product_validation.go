package middleware

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/flagsibh/mservices/product-api/data"
	"github.com/flagsibh/mservices/product-api/handlers"
	"github.com/flagsibh/mservices/product-api/utils"
	"github.com/go-playground/validator"
	hclog "github.com/hashicorp/go-hclog"
)

// Validation contains the validator to apply validation rules
type Validation struct {
	l        hclog.Logger
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
func NewValidation(l hclog.Logger) *Validation {
	validate := validator.New()
	validate.RegisterValidation("sku", skuValidator)
	return &Validation{l, validate}
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

// ProductValidationMiddleware is middleware for product validation
func (v *Validation) ProductValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := utils.FromJSON(prod, r.Body)
		if err != nil {
			v.l.Error("deserializing product", err)

			rw.WriteHeader(http.StatusBadRequest)
			utils.ToJSON(&data.ErrGenericError{Message: err.Error()}, rw)
			return
		}

		// validate the product
		errs := v.Validate(prod)
		if len(errs) != 0 {
			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			utils.ToJSON(&data.ErrGenericErrors{Messages: errs.Errors()}, rw)
			return
		}

		ctx := context.WithValue(r.Context(), handlers.KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
