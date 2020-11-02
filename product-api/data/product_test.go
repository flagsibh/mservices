package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name:  "flags",
		Price: 1.00,
		SKU:   "abc-fgj-lmn",
	}
	err := NewValidation().Validate(p)

	if err != nil {
		t.Fatal(err)
	}
}
