package docs

import (
	"github.com/flagsibh/mservices/product-api/data"
	"github.com/flagsibh/mservices/product-api/handlers"
)

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body data.ErrGenericError
}

// Validation errors defined as an array of strings
// swagger:response errorValidation
type errorValidationResponseWrapper struct {
	// Collection of the errors
	// in: body
	Body data.ErrGenericErrors
}

// Data structure representing a single product
// swagger:response productResponse
type productResponseWrapper struct {
	// Newly created product
	// in: body
	Body data.Product
}

// Response containing an array of products.
// swagger:response productsResponse
type productsResponseWrapper struct {
	// in:body
	Body handlers.ProductsResponse
}

// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// swagger:parameters deleteProduct
type idParameterWrapper struct {
	// The id of the product to delete
	// in:path
	// required: true
	ID int `json:"id"`
}
