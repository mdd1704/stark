CREATE TABLE IF NOT EXISTS user_details
(
    id CHAR(36) PRIMARY KEY COMMENT 'Identifier',
    device_token VARCHAR(255) COMMENT 'Device Token',
    device_os VARCHAR(25) COMMENT 'Device OS',
    avatar_url VARCHAR(255) COMMENT 'Avatar URL',
    avatar_path VARCHAR(255) COMMENT 'Avatar Path',
    source VARCHAR(50) COMMENT 'Source',
    oauth_id VARCHAR(255) UNIQUE COMMENT 'OAuth ID',
    id_card_url VARCHAR(255) COMMENT 'ID Card URL',
    id_card_path VARCHAR(255) COMMENT 'ID Card Path',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT 'Created At',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated At',
    CONSTRAINT detail_user_fk FOREIGN KEY (id) REFERENCES users (id) ON DELETE CASCADE
) COMMENT 'User Details' CHARSET=utf8;




