package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/flagsibh/mservices/product-api/data"
	"github.com/gorilla/mux"
)

// Products collection of products
type Products struct {
	l *log.Logger
}

// NewProducts creates a new product list
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// GetProducts get a list of products
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// AddProduct Creates a new product
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
}

// UpdateProduct updates a product
func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
	}

	p.l.Printf("Got ID: %d", id)
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	errnf := data.UpdateProduct(id, prod)
	if errnf == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
	} else {
		if errnf != nil {
			http.Error(rw, "Product not found", http.StatusInternalServerError)
		}
	}

}

// KeyProduct to identifiy the product in the context
type KeyProduct struct{}

//ProductValidation is middleware for product validation
func (p *Products) ProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		er := prod.FromJSON(r.Body)
		if er != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		// validate the product
		err := prod.Validate()
		if err != nil {
			http.Error(rw, fmt.Sprintf("Error validating the product: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
