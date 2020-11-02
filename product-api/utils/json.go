package utils

import (
	"encoding/json"
	"io"
)

// FromJSON reads JSON from the requests body
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}

// ToJSON convert data to JSON
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}
