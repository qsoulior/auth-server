package v1

import "net/http"

func Token(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// create >> signin -> access+refresh
	// update -> access+refresh
	// delete >> signout
}
