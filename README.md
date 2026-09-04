# IT Support API

A RESTful API for managing IT support tickets, built with Go and SQLite.

The API allows users to create, retrieve, update, and delete support tickets, with additional features for searching, filtering, pagination, sorting, authentication, and request logging.

## Features

* Ticket CRUD operations
* Ticket search
* Filtering by status and priority
* Pagination
* Sorting
* API key authentication
* Request logging
* SQLite database
* Environment-based configuration


## API Endpoints

All `/tickets` endpoints require an API key.

### Get all tickets

```text
GET /tickets
```

Returns a list of tickets.

Example:

```bash
curl -H "X-API-Key: your-hidden-key" http://localhost:8090/tickets
```

### Get a single ticket

```text
GET /tickets/{id}
```

Example:

```bash
curl -H "X-API-Key: your-hidden-key" http://localhost:8090/tickets/3
```

### Create a ticket

```text
POST /tickets
```

Example:

```bash
curl -X POST \
  -H "X-API-Key: your-hidden-key" \
  -H "Content-Type: application/json" \
  -d '{"title":"Keyboard not working","description":"User cannot type on the keyboard","priority":"high"}' \
  http://localhost:8090/tickets
```

### Update a ticket

```text
PATCH /tickets/{id}
```

Example:

```bash
curl -X PATCH \
  -H "X-API-Key: your-hidden-key" \
  -H "Content-Type: application/json" \
  -d '{"status":"in_progress"}' \
  http://localhost:8090/tickets/3
```

### Delete a ticket

```text
DELETE /tickets/{id}
```

Example:

```bash
curl -X DELETE \
  -H "X-API-Key: your-hidden-key" \
  http://localhost:8090/tickets/3
```
## Query Parameters

The `GET /tickets` endpoint supports filtering, searching, pagination, and sorting.

### Filter by status

```text
GET /tickets?status=open
```

Example:

```bash
curl -H "X-API-Key: your-hidden-key" \
  "http://localhost:8090/tickets?status=open"
```

### Filter by priority

```text
GET /tickets?priority=high
```

### Search tickets

```text
GET /tickets?search=wifi
```

Searches ticket information for the provided search term.

### Pagination

Use `page` and `limit` to control the number of results returned.

```text
GET /tickets?page=1&limit=10
```

### Sorting

Use `sort` to select the field and `order` to choose ascending or descending order.

```text
GET /tickets?sort=id&order=desc
```

Supported order values:

* `asc` — ascending
* `desc` — descending

### Combine parameters

Parameters can be combined in a single request.

Example:

```bash
curl -H "X-API-Key: your-hidden-key" \
  "http://localhost:8090/tickets?status=open&priority=high&sort=id&order=desc"
```

## Getting Started

### Prerequisites

Make sure you have the following installed:

* Go
* Git

### Clone the repository

```bash
git clone https://github.com/Billi070/it-support-api.git
cd it-support-api
```

### Install dependencies

```bash
go mod tidy
```

### Configure environment variables

Create a `.env` file in the project root:

```text
API_KEY=your-hidden-key
```

The `.env` file is intentionally excluded from Git.

### Run the API

```bash
go run .
```

The server will start on:

```text
http://localhost:8090
```

### Test the API

Send a request using your configured API key:

```bash
curl -H "X-API-Key: your-hidden-key" \
  http://localhost:8090/tickets
```

You should receive a JSON response containing the available tickets.

## Technology Stack

* **Go** — REST API and backend logic
* **SQLite** — Database
* **net/http** — HTTP server and routing
* **godotenv** — Environment configuration
* **Git/GitHub** — Version control

