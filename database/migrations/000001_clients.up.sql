CREATE TABLE IF NOT EXISTS clients
(
    id CHAR(36) PRIMARY KEY COMMENT 'Identifier',
    name VARCHAR(100) COMMENT 'Name',
    bearer_key VARCHAR(255) unique COMMENT 'Bearer Key',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT 'Created At',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated At'
) COMMENT 'Clients' CHARSET=utf8;