package asciiif

import (
	"encoding/json"
	"io"
)

// EncodeAll encodes the given ASCIIIF to json and writes it to w.
func EncodeAll(w io.Writer, a *ASCIIIF) error {
	encoder := json.NewEncoder(w)

	err := encoder.Encode(a)
	if err != nil {
		return err
	}
	return nil
}
