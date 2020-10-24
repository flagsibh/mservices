package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/flagsibh/mservices/product-api/data"
)

// Products collection of products
type Products struct {
	l *log.Logger
}

// NewProducts creates a new product list
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		p.updateProduct(rw, r)
		return
	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}
	data.AddProduct(prod)
}

func (p *Products) updateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT Product")

	regex := regexp.MustCompile(`/([0-9]+)`)
	res := regex.FindAllStringSubmatch(r.URL.Path, -1)

	if len(res) != 1 || len(res[0]) != 2 {
		http.Error(rw, "Invalid URL", http.StatusBadRequest)
	}

	id, err := strconv.Atoi(res[0][1])
	if err != nil {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
	}

	p.l.Printf("Got ID: %d", id)

	prod := &data.Product{}
	er := prod.FromJSON(r.Body)
	if er != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	errnf := data.UpdateProduct(id, prod)
	if errnf == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
	} else {
		if errnf != nil {
			http.Error(rw, "Product not found", http.StatusInternalServerError)
		}
	}

}
