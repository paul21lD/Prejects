package utils

import (
	"encoding/json"
	"net/http"
)

func ParseBodyStrict(w http.ResponseWriter, r *http.Request, dst any) error {
	// limita la dimensione del body (es. 1MB)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // errore se ci sono campi extra nel JSON

	if err := dec.Decode(dst); err != nil {
		return err // il caller può rispondere 400 Bad Request
	}
	// opzionale: assicurati che non ci siano dati extra dopo il JSON
	if dec.More() {
		return http.ErrBodyNotAllowed
	}
	return nil
}
