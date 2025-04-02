# Backend Project Structure

This project follows a clean and organized backend structure using Go.

## Folder Structure
```
Backend/
├── cmd/              # Entry points (main.go)
├── config/           # Configuration files
├── data/             # Data respitory (replacing database)
├── internal/         # Private app logic
│   ├── server/       # HTTP server setup
│   ├── handlers/     # API controllers (request handling)
│   ├── services/     # Business logic (data processing)
│   ├── models/       # Data models (structs)
│   ├── middleware/   # Middleware functions (auth, logging, etc.)
│   ├── utils/        # Utility functions (helpers)
|       ├── AI_Funcs  # AI functions (helpers)
|       ├── JD_Manage # JD adding functions
├── migrations/       # Database migration scripts
├── go.mod            # Go module file
├── Dockerfile        # Docker setup
├── README.md       # Documentation
```

## Getting Started

### 1. Run the API Server
```sh
go run cmd/main.go
```

### 1. API Endpoints
```
| Method | Endpoint  | Description      |
|--------|----------|-----------------|
| GET    | /user    | Fetch user data |



...

```

### 4. Docker Setup
To run with Docker:
```sh
docker build -t backend .
docker run -p 8080:8080 backend
```

---

