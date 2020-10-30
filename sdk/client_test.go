package sdk

import (
	"fmt"
	"testing"

	"github.com/flagsibh/mservices/sdk/client"
	"github.com/flagsibh/mservices/sdk/client/products"
)

func TestClient(t *testing.T) {
	c := client.NewHTTPClientWithConfig(nil,
		client.DefaultTransportConfig().WithHost("localhost:9090"))
	prods, err := c.Products.ListProducts(products.NewListProductsParams())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(prods)
}
