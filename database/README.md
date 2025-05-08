# Database Setup

This folder contains the shared SQLite database used by all three backend implementations.

## Initial Setup

To initialize the database with test users, run:

```bash
cd database
go run init.go
```

This will create a `login.db` file with the following test users:
1. Username: test, Password: test123
2. Username: admin, Password: admin123

## Database Structure

The database contains a single table:

### Users Table
- id: INTEGER PRIMARY KEY AUTOINCREMENT
- username: TEXT UNIQUE NOT NULL
- password: TEXT NOT NULL (bcrypt hashed)
- created_at: TIMESTAMP DEFAULT CURRENT_TIMESTAMP
