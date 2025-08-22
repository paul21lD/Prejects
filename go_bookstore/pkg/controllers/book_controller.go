package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/paul21ID/Prejects/go_bookstore/pkg/models"
	"github.com/paul21ID/Prejects/go_bookstore/pkg/utils"

	"github.com/gorilla/mux"
)

var NewBook models.Book

func GetBook(w http.ResponseWriter, r *http.Request) {
	books, err := models.GetAllBooks()
	if err != nil {
		http.Error(w, "errore DB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(books); err != nil {
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
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "libro non trovato", http.StatusNotFound)
		} else {
			http.Error(w, "errore DB: "+err.Error(), http.StatusInternalServerError)
		}
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

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookId := mux.Vars(r)["bookId"]
	ID, err := strconv.Atoi(bookId)
	if err != nil || ID <= 0 {
		http.Error(w, "id non valido", http.StatusBadRequest)
		return
	}

	book, err := models.DeleteBook(int64(ID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "libro non trovato", http.StatusNotFound)
		} else {
			http.Error(w, "errore DB: "+err.Error(), http.StatusInternalServerError)
		}
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

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// id
	bookId := mux.Vars(r)["bookId"]
	ID, err := strconv.Atoi(bookId)
	if err != nil || ID <= 0 {
		http.Error(w, "id non valido", http.StatusBadRequest)
		return
	}

	// leggi body con parsing + 400 error in caso
	update := &models.Book{}
	if err := utils.ParseBodyStrict(w, r, update); err != nil {
		http.Error(w, "JSON non valido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// carica esistente
	book, err := models.GetBookById(int64(ID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "libro non trovato", http.StatusNotFound)
		} else {
			http.Error(w, "errore DB: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if update.Name != "" {
		book.Name = update.Name
	}
	if update.Author != "" {
		book.Author = update.Author
	}
	if update.Publication != "" {
		book.Publication = update.Publication
	}

	if err := models.UpdateBook(book); err != nil {
		http.Error(w, "errore DB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(book)
}
