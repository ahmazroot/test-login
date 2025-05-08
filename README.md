# Login System Monorepo

This is a monorepo containing a Next.js frontend and three different backend implementations for a login system.

## Project Structure

- `frontend/` - Next.js frontend with shadcn/ui
- `backend-axum/` - Rust Axum backend
- `backend-rocket/` - Rust Rocket backend
- `backend-go/` - Go Fiber backend

## Getting Started

### Frontend
```bash
cd frontend
npm install
npm run dev
```

### Backend (Axum)
```bash
cd backend-axum
cargo run
```

### Backend (Rocket)
```bash
cd backend-rocket
cargo run
```

### Backend (Go Fiber)
```bash
cd backend-go
go run main.go
```

## Features

- User authentication with username and password
- SQLite database for storing user credentials
- Modern UI with shadcn/ui components
- Multiple backend implementations showcasing different frameworks
