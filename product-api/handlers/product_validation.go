package handlers

import (
	"context"
	"net/http"

	"github.com/flagsibh/mservices/product-api/data"
	"github.com/flagsibh/mservices/product-api/utils"
)

//ProductValidation is middleware for product validation
func (p *Products) ProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := utils.FromJSON(prod, r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)

			rw.WriteHeader(http.StatusBadRequest)
			utils.ToJSON(&data.ErrGenericError{Message: err.Error()}, rw)
			return
		}

		// validate the product
		errs := p.v.Validate(prod)
		if len(errs) != 0 {
			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			utils.ToJSON(&data.ErrGenericErrors{Messages: errs.Errors()}, rw)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
