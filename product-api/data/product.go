package data

import (
	"fmt"
	"time"
)

// Product is the main data unit
// swagger:model
type Product struct {
	// the id of the product.
	// required: true
	// min: 1
	ID int `json:"id"`
	// the name of the product
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`
	// the description of the product
	// required: false
	// max length: 10000
	Description string `json:"description"`
	// the price of the product
	// required: true
	// min: 0.01
	Price float32 `json:"price" validate:"gt=0"`
	// the unique product identification
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU       string `json:"sku" validate:"required,sku"`
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

// Products list of products
type Products []*Product

// GetProducts returns a list of products
func GetProducts() Products {
	return productList
}

// AddProduct adds a new product to the list
func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

// UpdateProduct updates a product on the list
func UpdateProduct(id int, p *Product) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	p.ID = id
	productList[i] = p
	return nil
}

// DeleteProduct deletes a product from the database
func DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

var productList = Products{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

// ErrProductNotFound error to signal product does not exists
var ErrProductNotFound = fmt.Errorf("Product not found")

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

// FindProduct returns a single product which matches the id.
// If the product is not found, this functions returns a ErrProductNotFound error.
func FindProduct(id int) (*Product, error) {
	i := findIndexByProductID(id)

	if i == -1 {
		return nil, ErrProductNotFound
	}
	return productList[i], nil
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}
