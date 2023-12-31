DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id VARCHAR(32) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);