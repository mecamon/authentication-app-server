package helpers

import "net/http"

func ResGenerator(w http.ResponseWriter, code int, output []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(output)
}
