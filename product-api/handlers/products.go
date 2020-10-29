package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/flagsibh/mservices/product-api/data"
	"github.com/flagsibh/mservices/product-api/utils"
	"github.com/gorilla/mux"
)

// KeyProduct to identifiy the product in the context
type KeyProduct struct{}

// Products collection of products
type Products struct {
	l *log.Logger
	v *data.Validation
}

// NewProducts creates a new product list
func NewProducts(l *log.Logger, v *data.Validation) *Products {
	return &Products{l, v}
}

// GetProducts get a list of products
// swagger:route GET / products listProducts
// Returns a list of products.
// Responses:
//	200: productsResponse
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	lp := data.GetProducts()
	err := utils.ToJSON(lp, rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// GetProduct get a product from the list
// swagger:route GET /{id} products findProduct
// Retuns a single product from the list.
// Responses:
//	200: productResponse
// 	404: errorResponse
func (p *Products) GetProduct(rw http.ResponseWriter, r *http.Request) {
	id := p.getProductID(r)

	prod, err := data.FindProduct(id)

	switch err {
	case nil:
	case data.ErrProductNotFound:
		p.l.Println("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusNotFound)
		utils.ToJSON(&data.ErrGenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Println("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusInternalServerError)
		utils.ToJSON(&data.ErrGenericError{Message: err.Error()}, rw)
		return
	}

	err = utils.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Println("[ERROR] serializing product", err)
	}
}

// AddProduct Creates a new product
// swagger:route POST / products createProduct
// Creates a new product.
// Responses:
//	201: productResponse
// 	422: errorValidationResponse
//	501: errorResponse
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
}

// UpdateProduct updates a product
// swagger:route PUT / products updateProduct
// Updates an existing product.
// Responses:
// 	200: noContentResponse
// 	404: errorResponse
// 	422: errorValidationResponse
func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")

	id := p.getProductID(r)

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err := data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		rw.WriteHeader(http.StatusNotFound)
		utils.ToJSON(&data.ErrGenericError{Message: "Product not found in database"}, rw)
		return
	}

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)
}

// DeleteProduct deletes a product from the list
// swagger:route DELETE /{id} products deleteProduct
// Deletes a product.
// Responses:
//	204: noContentResponse
//	404: errorResponse
//	501: errorResponse
func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle DELETE Product")

	id := p.getProductID(r)

	p.l.Printf("Got ID: %d", id)

	err := data.DeleteProduct(id)

	if err == data.ErrProductNotFound {
		rw.WriteHeader(http.StatusNotFound)
		utils.ToJSON(&data.ErrGenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.ToJSON(&data.ErrGenericError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (p *Products) getProductID(r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// http.Error(rw, "Invalid ID", http.StatusBadRequest)
		panic(err)
	}

	p.l.Printf("Got ID: %d", id)

	return id
}
