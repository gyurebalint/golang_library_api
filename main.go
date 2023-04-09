package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/mux"
)

type book struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	Has_read    bool   `json:"has_read"`
}

type response struct {
	HttpStatusCode int         `json:"httpStatusCode"`
	SuccessMessage string      `json:"successMessage"`
	Body           interface{} `json:"body"`
}

var books = []book{
	{Id: "1", Title: "C# Garbage Collection", Author: "Some C# Guru", Description: "Long book about garbage collection", Genre: "Technical", Has_read: false},
	{Id: "2", Title: "Then there were none", Author: "Agatha Christie", Description: "An eerie whodunit", Genre: "Entertainment", Has_read: true},
	{Id: "3", Title: "System Design interview", Author: "alex Xu", Description: "A good book on system design interviews", Genre: "Technical", Has_read: true},
}

func respondWithJSON(writer http.ResponseWriter, httpStatusCode int, successMessage string, payload interface{}) {
	writer.WriteHeader(200)
	p := response{
		HttpStatusCode: httpStatusCode,
		SuccessMessage: successMessage,
		Body:           payload,
	}

	writer.Header().Set("Content-Type", "application/json")
	resp, err := json.MarshalIndent(p, "  ", "  ")
	if err != nil {
		log.Fatal()
	}
	writer.Write(resp)
}

// getAlbums responds with the list of all albums as JSON.
func getBooks(writer http.ResponseWriter, request *http.Request) {

	respondWithJSON(writer, http.StatusOK, "SUCCESS", books)
}

func getBookByID(writer http.ResponseWriter, request *http.Request) {
	id := path.Base(request.URL.Path)

	if id != "" {
		for _, b := range books {
			if b.Id == id {
				respondWithJSON(writer, http.StatusOK, "SUCCESS", []book{b})
				return
			}
		}
	}
	respondWithJSON(writer, http.StatusNoContent, "NOT FOUND", map[string]string{"message": "book not found"})

}

func updateBook(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	var b book
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&b); err != nil {
		respondWithJSON(writer, http.StatusBadRequest, "Invalid resquest payload", nil)
		return
	}

	defer request.Body.Close()
	b.Id = vars["id"]
	if b.Id == "" {
		respondWithJSON(writer, http.StatusNoContent, "NOT FOUND", map[string]string{"message": "Id in URL not found"})
		return
	}

	var index int = -1
	for i, book := range books {
		if book.Id == b.Id {
			index = i
			break
		}
	}
	if index == -1 {
		respondWithJSON(writer, http.StatusNoContent, "NOT FOUND", map[string]string{"message": "No such book with that id in the database"})
		return
	}
	books[index].Author = b.Author
	books[index].Description = b.Description
	books[index].Genre = b.Genre
	books[index].Has_read = b.Has_read
	books[index].Title = b.Title

	respondWithJSON(writer, http.StatusOK, "SUCCESS", books[index])
}

func addBook(writer http.ResponseWriter, request *http.Request) {
	var b book
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&b); err != nil {
		respondWithJSON(writer, http.StatusBadRequest, "Invalid resquest payload", nil)
		return
	}
	defer request.Body.Close()

	var maxId int
	for _, b := range books {
		bookId, err := strconv.Atoi(b.Id)
		if err != nil {
			log.Fatal()
		}
		if bookId > maxId {
			maxId = bookId
		}
	}
	b.Id = strconv.Itoa(maxId + 1)
	books = append(books, b)
	respondWithJSON(writer, http.StatusOK, "SUCCESSFULLY ADDED", b)
}

func deleteBook(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id := vars["id"]
	if id == "" {
		respondWithJSON(writer, http.StatusNoContent, "NOT FOUND", map[string]string{"message": "Id in URL not found"})
		return
	}

	var index int = -1
	for i, book := range books {
		if book.Id == id {
			index = i
			break
		}
	}

	if index == -1 {
		respondWithJSON(writer, http.StatusNoContent, "NOT FOUND", map[string]string{"message": "No such book with that id in the database"})
		return
	}

	books = append(books[:index], books[index+1:]...)
	respondWithJSON(writer, http.StatusOK, "SUCCESSFULLY DELETED", books)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", getBooks).Methods("GET")
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBookByID).Methods("GET")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books", addBook).Methods("POST")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	srv := &http.Server{
		Handler: r,
		Addr:    ":8080",
	}
	log.Fatal(srv.ListenAndServe())
}
