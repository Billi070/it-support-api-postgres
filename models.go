package main

type Ticket struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int  `json:"priority"`
	Status      string `json:"status"`
}