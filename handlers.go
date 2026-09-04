package main

import (
	sql "database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//var tickets []Ticket

func ticketsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tickets" {
		switch r.Method {
		case http.MethodGet:
			getTicket(w, r)
		case http.MethodPatch:
			updateTicket(w, r)
		case http.MethodDelete:
			deleteTicket(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	switch r.Method {
	case http.MethodPost:
		createTicket(w, r)

	case http.MethodGet:
		getTickets(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}

func validateTicket(ticket Ticket) error {
	if strings.TrimSpace(ticket.Title) == "" {
		return fmt.Errorf("title is required")
	}

	if strings.TrimSpace(ticket.Description) == "" {
		return fmt.Errorf("description is required")
	}

	// switch ticket.Priority {
	// case "low", "medium", "high":
	// 	return nil
	// default:
	// 	return fmt.Errorf("priority must be low, medium, or high")
	// }

	return validatePriority(ticket.Priority)

}

func validateStatus(status string) error {
	switch status {
	case "open", "in_progress", "resolved", "closed":
		return nil
	default:
		return fmt.Errorf("invalid status")
	}
}

func createTicket(w http.ResponseWriter, r *http.Request) {
	var ticket Ticket

	err := json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validateTicket(ticket); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ticket.Status = "open"

	result, err := db.Exec(
		"INSERT INTO tickets (title, description, priority, status) VALUES (?, ?, ?, ?)",
		ticket.Title,
		ticket.Description,
		ticket.Priority,
		ticket.Status,
	)
	if err != nil {
		http.Error(w, "Failed to create ticket", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get ticket ID", http.StatusInternalServerError)
		return
	}

	ticket.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticket)
}

func getTickets(w http.ResponseWriter, r *http.Request) {

	// 	rows, err := db.Query(`
	// 	SELECT id, title, description, priority, status
	// 	FROM tickets
	// `)

	// if err != nil {
	// 	http.Error(w, "Failed to get tickets", http.StatusInternalServerError)
	// 	return
	// }

	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")
	search := r.URL.Query().Get("search")
	// page := r.URL.Query().Get("page")
	// limit := r.URL.Query().Get("limit")

	// pageString := r.URL.Query().Get("page")
	// limitString := r.URL.Query().Get("limit")

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	if sort == "" {
		sort = "id"
	}

	if order == "" {
		order = "asc"
	}

	switch sort {
	case "id", "title", "priority", "status":
		// valid
	default:
		http.Error(w, "Invalid sort field", http.StatusBadRequest)
		return
	}

	if order != "asc" && order != "desc" {
		http.Error(w, "Invalid sort order", http.StatusBadRequest)
		return
	}

	// page := 1
	// limit := 10
	// var err error

	// if pageString != "" {

	// 	page, err = strconv.Atoi(pageString)
	// 	if err != nil || page < 1 {
	// 		http.Error(w, "Invalid page", http.StatusBadRequest)
	// 		return
	// 	}
	// }

	// if limitString != "" {

	// 	limit, err = strconv.Atoi(limitString)
	// 	if err != nil || limit < 1 || limit > 100{
	// 		http.Error(w, "Limit must be between 1 and 100", http.StatusBadRequest)
	// 		return
	// 	}
	// }

	// offset := (page - 1) * limit

	page, limit, err := getPagination(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset := (page - 1) * limit

	if status != "" {
		if err := validateStatus(status); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if priority != "" {
		if err := validatePriority(priority); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	query := `
	SELECT id, title, description, priority, status
	FROM tickets
	WHERE 1=1
`

	args := []any{}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	if priority != "" {
		query += " AND priority = ?"
		args = append(args, priority)
	}

	if search != "" {
		query += " AND (title LIKE ? OR description LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	countQuery := "SELECT COUNT(*) FROM tickets WHERE 1=1"
	countArgs := []any{}

	if status != "" {
		countQuery += " AND status = ?"
		countArgs = append(countArgs, status)
	}

	if priority != "" {
		countQuery += " AND priority = ?"
		countArgs = append(countArgs, priority)
	}

	if search != "" {
		countQuery += " AND (title LIKE ? OR description LIKE ?)"
		searchTerm := "%" + search + "%"
		countArgs = append(countArgs, searchTerm, searchTerm)
	}

	query += " ORDER BY " + sort + " " + strings.ToUpper(order)
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "Failed to get tickets", http.StatusInternalServerError)
		return
	}

	var total int

	err = db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		http.Error(w, "Failed to count tickets", http.StatusInternalServerError)
		return
	}

	totalPages := (total + limit - 1) / limit

	defer rows.Close()

	tickets := []Ticket{}

	for rows.Next() {
		var ticket Ticket

		err := rows.Scan(
			&ticket.ID,
			&ticket.Title,
			&ticket.Description,
			&ticket.Priority,
			&ticket.Status,
		)
		if err != nil {
			http.Error(w, "Failed to read ticket", http.StatusInternalServerError)
			return
		}

		tickets = append(tickets, ticket)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to read tickets", http.StatusInternalServerError)
		return
	}

	response := struct {
		Tickets    []Ticket `json:"tickets"`
		Page       int      `json:"page"`
		Limit      int      `json:"limit"`
		Total      int      `json:"total"`
		TotalPages int      `json:"total_pages"`
	}{
		Tickets:    tickets,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKey := r.Header.Get("X-API-Key")
		expectedKey := os.Getenv("API_KEY")

		// 		fmt.Println("Expected:", expectedKey)
		// fmt.Println("Received:", apiKey)

		if apiKey != expectedKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// func (rw *responseWriter) WriteHeader(statusCode int) {
// 	rw.statusCode = statusCode
// 	rw.ResponseWriter.WriteHeader(statusCode)
// }

func (rw *responseWriter) Write(data []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}

	return rw.ResponseWriter.Write(data)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if rw.statusCode != 0 {
		return
	}

	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		//fmt.Println(r.Method, r.URL.Path)

		rw := &responseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(rw, r)
		duration := time.Since(start)
		fmt.Println(r.Method, r.URL.Path, rw.statusCode, duration)
	})
}

func getPagination(r *http.Request) (int, int, error) {
	pageString := r.URL.Query().Get("page")
	limitString := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	var err error

	if pageString != "" {
		page, err = strconv.Atoi(pageString)
		if err != nil || page < 1 {
			return 0, 0, fmt.Errorf("invalid page")
		}
	}

	if limitString != "" {
		limit, err = strconv.Atoi(limitString)
		if err != nil || limit < 1 || limit > 100 {
			return 0, 0, fmt.Errorf("limit must be between 1 and 100")
		}
	}

	return page, limit, nil
}

func getTicket(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/tickets/")

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}
	var ticket Ticket

	err = db.QueryRow(`
		SELECT id, title, description, priority, status
		FROM tickets
		WHERE id = ?
	`, id).Scan(
		&ticket.ID,
		&ticket.Title,
		&ticket.Description,
		&ticket.Priority,
		&ticket.Status,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Failed to get ticket", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticket)
}

func updateTicket(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/tickets/")

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	var update struct {
		Status   string `json:"status"`
		Priority string `json:"priority"`
	}

	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//		if err := validateStatus(update.Status); err != nil {
	//		http.Error(w, err.Error(), http.StatusBadRequest)
	//		return
	//	}
	if update.Status != "" {
		if err := validateStatus(update.Status); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if update.Priority != "" {
		if err := validatePriority(update.Priority); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// result, err := db.Exec(
	// 	"UPDATE tickets SET status = ? WHERE id = ?",
	// 	update.Status,
	// 	id,
	// )

	updates := []string{}
	args := []any{}

	if update.Status != "" {
		updates = append(updates, "status = ?")
		args = append(args, update.Status)
	}

	if update.Priority != "" {
		updates = append(updates, "priority = ?")
		args = append(args, update.Priority)
	}

	if len(updates) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	args = append(args, id)

	query := "UPDATE tickets SET " + strings.Join(updates, ", ") + " WHERE id = ?"

	result, err := db.Exec(query, args...)
	if err != nil {
		http.Error(w, "Failed to update ticket", http.StatusInternalServerError)
		return
	}

	// if err != nil {
	// 	http.Error(w, "Failed to update ticket", http.StatusInternalServerError)
	// 	return
	// }

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check update", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	var ticket Ticket

	err = db.QueryRow(`
		SELECT id, title, description, priority, status
		FROM tickets
		WHERE id = ?
	`, id).Scan(
		&ticket.ID,
		&ticket.Title,
		&ticket.Description,
		&ticket.Priority,
		&ticket.Status,
	)

	if err != nil {
		http.Error(w, "Failed to get updated ticket", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticket)
}

func validatePriority(priority string) error {
	switch priority {
	case "low", "medium", "high":
		return nil
	default:
		return fmt.Errorf("priority must be low, medium, or high")
	}
}

func deleteTicket(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/tickets/")

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(
		"DELETE FROM tickets WHERE id = ?",
		id,
	)
	if err != nil {
		http.Error(w, "Failed to delete ticket", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check deletion", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
