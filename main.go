package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	initDB()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env")
		return
	}

	//http.HandleFunc("/tickets", ticketsHandler)
	//http.Handle("/tickets", authMiddleware(http.HandlerFunc(ticketsHandler)))
	//http.Handle("/tickets/", authMiddleware(http.HandlerFunc(ticketsHandler)))
	//http.HandleFunc("/tickets/", ticketsHandler)

	http.Handle("/tickets", loggingMiddleware(authMiddleware(http.HandlerFunc(ticketsHandler))))
	http.Handle("/tickets/", loggingMiddleware(authMiddleware(http.HandlerFunc(ticketsHandler))))
	fmt.Println("Server is running on http://localhost:8090")
	err = http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
