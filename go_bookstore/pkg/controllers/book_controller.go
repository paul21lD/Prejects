package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/paul21ID/Prejects/go_bookstore/pkg/models"
	"github.com/paul21ID/Prejects/go_bookstore/pkg/utils"

	"github.com/gorilla/mux"
)

var NewBook models.Book

func getBook(w http.ResponseWriter, r *http.Request) {
	books, err := models.GetAllBooks()
	if err != nil {
		http.Error(w, "errore DB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newBooks); err != nil {
		// qui sei già “in scrittura”: logga e chiudi
		// (se vuoi essere ultra-sicuro, potresti encodare prima in un buffer)
		// log.Printf("encode error: %v", err)
		return
	}
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	bookId := mux.Vars(r)["bookId"]
	ID, err := strconv.Atoi(bookId)
	if err != nil || ID <= 0 {
		http.Error(w, "id non valido", http.StatusBadRequest)
		return
	}

	book, err := models.GetBookById(int64(ID))
	if err != nil {
		// se il tuo models distingue not-found dagli altri errori, mappa a 404/500
		http.Error(w, "libro non trovato", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(book); err != nil {
		// qui sei già “in scrittura”: logga e chiudi
		// (se vuoi essere ultra-sicuro, potresti encodare prima in un buffer)
		// log.Printf("encode error: %v", err)
		return
	}
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse del body: se JSON malformato → 400
	var in models.Book
	if err := utils.ParseBodyStrict(w, r, &in); err != nil {
		http.Error(w, "JSON non valido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validazione minima: se mancano campi obbligatori → 400
	if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Author) == "" {
		http.Error(w, "name e author sono obbligatori", http.StatusBadRequest)
		return
	}

	// Persistenza
	created := in.CreateBook()

	// Response 201 + JSON
	w.Header().Set("Content-Type", "application/json")
	if created.ID != 0 {
		w.Header().Set("Location", fmt.Sprintf("/book/%d", created.ID))
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(created)
}
