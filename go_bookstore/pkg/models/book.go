package models

import (
	"github.com/jinzhu/gorm"
	"github.com/paul21ID/Prejects/go_bookstore/pkg/config"
)

var db *gorm.DB

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Book{})
}

func (b *Book) CreateBook() *Book {
	db.NewRecord(b)
	db.Create(&b)
	return b
}

func GetAllBooks() ([]Book, error) {
	var books []Book
	tx := db.Find(&books)  // tx è *gorm.DB
	return books, tx.Error // restituisci l'eventuale errore
}

func GetBookById(Id int64) (*Book, error) {
	var book Book
	if err := db.First(&book, Id).Error; err != nil { // usa .Error dal *gorm.DB
		return nil, err
	}
	return &book, nil
}

func DeleteBook(ID int64) (Book, error) {
	var b Book
	// 1) carica il record
	if err := db.First(&b, ID).Error; err != nil {
		return Book{}, err // gorm.ErrRecordNotFound o altro
	}
	// 2) elimina (soft delete se usi gorm.Model)
	if err := db.Delete(&b).Error; err != nil {
		return Book{}, err
	}
	return b, nil
}

func UpdateBook(b *Book) error {
	return db.Save(b).Error
}
