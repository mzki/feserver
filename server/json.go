package server

import (
	"encoding/json"
	"io"

	"github.com/mzki/feserver/src"
)

// it represents json response returned from the server.
type JSONResponse struct {
	src.Response
	Error string `json:"error"`
}

func (res *JSONResponse) write(w io.Writer) error {
	return writeJSON(w, res)
}

func writeJSON(w io.Writer, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}
