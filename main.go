package main

import (
	"fmt"
	"groupie-tracker/handlers"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/", handlers.Mainpage)
	http.HandleFunc("/artist", handlers.Artistpage)

	fmt.Println("Server is running at http://localhost:8080")

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}
}
