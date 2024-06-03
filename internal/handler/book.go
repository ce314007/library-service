package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type BookHandler struct {
	db *sql.DB
}

type Book struct {
	ID            int     `json:"id" db:"id"`
	Title         string  `json:"title" db:"title"`
	Price         float64 `json:"price" db:"price"`
	PublishedDate string  `json:"publishedDate" db:"published_date"`
}

func NewBookHandler(db *sql.DB) *BookHandler {
	return &BookHandler{
		db: db,
	}
}

func (h *BookHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := h.db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "Database connection failed"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

func (h *BookHandler) GetAllBooksHandler(w http.ResponseWriter, r *http.Request) {
	var books []Book
	rows, err := h.db.Query("SELECT id, title, price, published_date FROM books")
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Price, &book.PublishedDate); err != nil {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]Book{"books": books})
}

func (h *BookHandler) AddNewBookHandler(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		fmt.Println(err)
		http.Error(w, "Bad request.", http.StatusBadRequest)
		return
	}

	err := h.db.QueryRow(
		"INSERT INTO books (title, price, published_date) VALUES ($1, $2, $3) RETURNING id",
		book.Title, book.Price, book.PublishedDate,
	).Scan(&book.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":            book.ID,
		"title":         book.Title,
		"price":         book.Price,
		"publishedDate": book.PublishedDate,
		"message":       "Book successfully added to the library.",
	})
}

func (h *BookHandler) GetBookByIDHandler(w http.ResponseWriter, r *http.Request, id string) {
	var book Book
	err := h.db.QueryRow("SELECT id, title, price, published_date FROM books WHERE id = $1", id).Scan(&book.ID, &book.Title, &book.Price, &book.PublishedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Book not found.", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) UpdateBookTitleHandler(w http.ResponseWriter, r *http.Request, id string) {
	var request struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Bad request.", http.StatusBadRequest)
		return
	}

	_, err := h.db.Exec("UPDATE books SET title = $1 WHERE id = $2", request.Title, id)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"title":   request.Title,
		"message": "Book title successfully updated.",
	})
}

func (h *BookHandler) DeleteBookByIDHandler(w http.ResponseWriter, r *http.Request, id string) {
	_, err := h.db.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Book successfully deleted.",
	})
}
