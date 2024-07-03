package jsonlib

import (
	"encoding/json"
	"io"
	"strings"
)

// UnmarshalReader reads from reader and decodes the JSON-encoded data into value.
func UnmarshalReader(reader io.Reader, value interface{}) error {
	d := json.NewDecoder(reader)
	for {
		if err := d.Decode(value); err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}
}

// MarshalString writes the JSON encodable value to a string.
func MarshalString(value interface{}) (string, error) {
	writer := &strings.Builder{}
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		return "", err
	}
	return writer.String(), nil
}
