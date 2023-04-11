package v1

import (
	"encoding/json"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	e := json.NewEncoder(w)
	e.Encode(map[string]string{
		"error": "Not Found",
	})
}
