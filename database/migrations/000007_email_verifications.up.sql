CREATE TABLE IF NOT EXISTS email_verifications
(
    email VARCHAR(100) COMMENT 'Email',
    token VARCHAR(255) unique COMMENT 'Token',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT 'Created At',
    CONSTRAINT verification_email_fk FOREIGN KEY (email) REFERENCES users (email) ON DELETE CASCADE
) COMMENT 'Email Verifications' CHARSET=utf8;




