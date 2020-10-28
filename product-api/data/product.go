package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
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

// Validate executes all validation definition/tags on a product
func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", skuValidator)
	return validate.Struct(p)
}

func skuValidator(fl validator.FieldLevel) bool {
	// SKU is abs-xsc-jhd
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return (len(matches) == 1)
}

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

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrProductNotFound
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
