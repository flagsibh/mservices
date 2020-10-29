package docs

import (
	"github.com/flagsibh/mservices/product-api/data"
)

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body data.ErrGenericError
}

// Validation errors defined as an array of strings
// swagger:response errorValidationResponse
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
	Body []data.Product
}

// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// swagger:parameters findProduct
type findProductParameterIDWrapper struct {
	// The id of the product to retrieve
	// in:path
	// required: true
	ID int `json:"id"`
}

// swagger:parameters deleteProduct
type deleteProductParameterIDWrapper struct {
	// The id of the product to delete
	// in:path
	// required: true
	ID int `json:"id"`
}
