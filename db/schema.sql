-- tickets table
CREATE TABLE tickets (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    status SmallINT NOT NULL DEFAULT 1,
    priority VARCHAR(50) NOT NULL,
    created_by BIGINT NOT NULL,
    assigned_to BIGINT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- users table
CREATE TABLE users (
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    full_name  VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    role       ENUM('admin','agent','customer') DEFAULT 'customer',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_tickets_created_by ON tickets(created_by);
