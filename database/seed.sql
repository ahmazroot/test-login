-- Create users table if not exists
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert test user if not exists
INSERT OR IGNORE INTO users (username, password) 
VALUES 
    -- Password: test123 (bcrypt hashed)
    ('test', '$2a$10$YourHashedPasswordHere'),
    -- Password: admin123 (bcrypt hashed)
    ('admin', '$2a$10$YourHashedPasswordHere');
