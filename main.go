package main

import (
	"encoding/json"
	"log"
	"net/http"
	"math/rand"
	"strconv"
	"fmt"
	"github.com/gorilla/mux"
)

// book struct and author struct
// struct just like a class
type Book struct {
	ID		string	`json:"id"`
	Isbn	string 	`json:"isbn"`
	Title	string 	`json:"title"`
	Author	*Author `json:"author"`
}

type Author struct {
	Firstname	string `json:"firstname"`
	LastName	string `json:"lastname"`
}


// init books var as a slice Book struct
// slice just like array, with dynamic length
var books []Book
func initData(w http.ResponseWriter, r *http.Request) {
	// create some mock data
	book1 := Book{ID: "1", Isbn: "11", Title: "Book One", Author: &Author{Firstname: "JP", LastName: "S"}}
	book2 := Book{ID: "2", Isbn: "22", Title: "Book Two", Author: &Author{Firstname: "JP", LastName: "S"}}
	book3 := Book{ID: "3", Isbn: "33", Title: "Book Three", Author: &Author{Firstname: "JP", LastName: "S"}}
	books = append(books, book1)
	books = append(books, book2)
	books = append(books, book3)
}


// get all books
func getBooks(w http.ResponseWriter, r *http.Request) {

	// set content-type of header
	w.Header().Set("Content-Type", "application/json")
	
	// encode books to JSON object and send to front-end by response writer
	json.NewEncoder(w).Encode(books)
}


// get single book
func getBook(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	// get params from the URL
	params := mux.Vars(r) 

	// loop through books and find with id
	// item is the iterator
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	// if no book was found in database, return nil Book struct
	json.NewEncoder(w).Encode(&Book{})
}


// create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	// create a book var, this book is from request data
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	// generate random book ID and Isbn
	book.ID = strconv.Itoa(rand.Intn(1000))
	book.Isbn = strconv.Itoa(rand.Intn(1000000))

	// add book to the database
	books = append(books, book)

	// return the book just created
	json.NewEncoder(w).Encode(book)
}


// update the book
func updateBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// get parameter from the URL
	params := mux.Vars(r)

	// for loop takes two parameters, index and item iterator
	for index, item := range books {
		if item.ID == params["id"] {

			// delete the book that need to update
			books = append(books[:index], books[index + 1:]...)

			// create a new book with updated info
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = item.ID
			book.Isbn = item.Isbn
			books = append(books, book)

			// return the book just created
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}


// delete a book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")

	// get parameter from the URL
	params := mux.Vars(r)

	// for loop takes two parameters, index and item iterator
	for index, item := range books {
		if item.ID == params["id"] {
			// use this to delete the item with index index
			books = append(books[:index], books[index + 1:]...)
			break
		}
	}

	// return all books after delete the book
	json.NewEncoder(w).Encode(books)
}


// main function
func main() {
	fmt.Println("Running!\n")

	// init router
	r := mux.NewRouter();

	// create route handlers and endpoints
	// use /param={id} to parse parameter to endpoint
	// at the end point, use params["id"] to get parameter in the request
	r.HandleFunc("/initdata", initData).Methods("GET")
	r.HandleFunc("/getbooks", getBooks).Methods("GET")
	r.HandleFunc("/getbook/{id}", getBook).Methods("GET")
	r.HandleFunc("/createbook", createBook).Methods("POST")
	r.HandleFunc("/updatebook/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/deletebook/{id}", deleteBook).Methods("DELETE")

	// start server, log if error occurs
	log.Fatal(http.ListenAndServe(":3000", r))
}
