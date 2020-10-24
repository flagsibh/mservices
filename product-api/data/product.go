package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Product is the main data unit
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// Products list of products
type Products []*Product

// FromJSON reads JSON from the requests body
func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

// ToJSON convert data to JSON
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

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
	_, i, err := findProduct(id)
	if err != nil {
		return err
	}

	p.ID = id
	productList[i] = p
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

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrProductNotFound
}
