-- 001_create_users_and_admin.sql

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- bcrypt hash for password "admin123"
-- You should generate your own secure password and hash!
INSERT INTO users (username, password_hash, role)
VALUES (
    'admin',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZag0u6b1l6Z1uE6Q7rQ9Qw1Q9Qw1Q', -- bcrypt hash for "admin123"
    'admin'
)
ON CONFLICT (username) DO NOTHING;
