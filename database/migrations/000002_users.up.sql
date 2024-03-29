CREATE TABLE IF NOT EXISTS users
(
    id CHAR(36) PRIMARY KEY COMMENT 'Identifier',
    name VARCHAR(100) COMMENT 'Name',
    email VARCHAR(100) UNIQUE COMMENT 'Email',
    username VARCHAR(25) UNIQUE COMMENT 'Username',
    contact VARCHAR(25) COMMENT 'Contact',
    password VARCHAR(255) COMMENT 'Password',
    email_verified_at TIMESTAMP DEFAULT NULL COMMENT 'Email Verified At',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT 'Created At',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated At'
) COMMENT 'Users' CHARSET=utf8;




