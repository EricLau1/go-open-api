CREATE DATABASE IF NOT EXISTS sandbox;

CREATE USER 'guest'@'%' IDENTIFIED BY 'guest';

GRANT ALL PRIVILEGES ON sandbox.* TO 'guest'@'%';

FLUSH PRIVILEGES;

USE sandbox;

CREATE TABLE IF NOT EXISTS users (
    id CHAR(36) PRIMARY KEY,
    email VARCHAR(50) UNIQUE NOT NULL,
    password CHAR(60) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP()
);