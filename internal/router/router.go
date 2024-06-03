package router

import (
	"database/sql"
	"gojek/library/internal/handler"
	"net/http"
)

func NewRouter(db *sql.DB) http.Handler {
	bookHandler := handler.NewBookHandler(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", handler.Ping)
	mux.HandleFunc("/healthz", bookHandler.HealthCheckHandler)
	mux.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bookHandler.GetAllBooksHandler(w, r)
		case http.MethodPost:
			bookHandler.AddNewBookHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		}
	})
	mux.HandleFunc("/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/books/"):]
		switch r.Method {
		case http.MethodGet:
			bookHandler.GetBookByIDHandler(w, r, id)
		case http.MethodPut:
			bookHandler.UpdateBookTitleHandler(w, r, id)
		case http.MethodDelete:
			bookHandler.DeleteBookByIDHandler(w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		}
	})

	return mux
}
