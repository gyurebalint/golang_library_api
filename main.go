package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	*Id          *int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	Has_read    bool   `json:"has_read"`
}

func (b book) toJson() {

	jsonBook, err := json.Marshal(b)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonBook))
}

// book slice to seed record album data.
var books = []book{
	{Id: 1, Title: "C# Garbage Collection", Author: "Some C# Guru", Description: "Long book about garbage collection", Genre: "Technical", Has_read: false},
	{Id: 2, Title: "Then there were none", Author: "Agatha Christie", Description: "An eerie whodunit", Genre: "Entertainment", Has_read: true},
	{Id: 3, Title: "System Design interview", Author: "alex Xu", Description: "A good book on system design interviews", Genre: "Technical", Has_read: true},
}

// getAlbums responds with the list of all albums as JSON.
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func addBook(c *gin.Context) {
	var newBook book
	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

/*
	"func getBook(id int) book {
		var book foundIt

		return
	}"
*/


func main() {
	// var b = book{Id: 4, Title: "System Design interview", Author: "alex Xu", Description: "A good book on system design interviews", Genre: "Technical", Has_read: true}
	// b.toJson()
	router := gin.Default()
	router.GET("/books", getBooks)
	//router.GET("/book/{id}", getBook)
	router.POST("/book", addBook)
	router.Run("localhost:8080")
}
