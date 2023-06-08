CREATE TABLE IF NOT EXISTS user_locations
(
    id CHAR(36) PRIMARY KEY COMMENT 'Identifier',
    province_id VARCHAR(50) COMMENT 'Province ID',
    regency_id VARCHAR(50) COMMENT 'Regency ID',
    district_id VARCHAR(50) COMMENT 'District ID',
    village_id VARCHAR(50) COMMENT 'Village ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT 'Created At',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated At',
    CONSTRAINT location_user_fk FOREIGN KEY (id) REFERENCES users (id) ON DELETE CASCADE
) COMMENT 'User Locations' CHARSET=utf8;




