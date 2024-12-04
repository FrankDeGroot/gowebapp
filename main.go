package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type Book struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}

func main() {

	InitDb()
	defer Pool.Close()

	fmt.Println("Starting web server")
	r := http.NewServeMux()
	r.HandleFunc("GET /v1/books/{id}", getBookHandler)
	http.ListenAndServe(":8000", r)
}

func getBookHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	switch err {
	case nil:
		break
	case strconv.ErrSyntax:
	case strconv.ErrRange:
		http.Error(w, "id not an integer", http.StatusBadRequest)
	default:
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Fprint(os.Stderr, err)
	}

	var title string
	err = Pool.QueryRow(context.Background(), "select title from book where id=$1", id).Scan(&title)
	switch err {
	case nil:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Book{id, title})
	case pgx.ErrNoRows:
		http.Error(w, "id not found", http.StatusNotFound)
	default:
		http.Error(w, "Error", http.StatusInternalServerError)
		fmt.Fprint(os.Stderr, err)
	}
}
