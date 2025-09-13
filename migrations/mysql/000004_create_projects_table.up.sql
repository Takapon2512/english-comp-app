CREATE TABLE IF NOT EXISTS projects (
    id CHAR(26) NOT NULL,
    user_id CHAR(26) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    PRIMARY KEY (id),
    INDEX user_id_idx (user_id),
    INDEX created_at_idx (created_at),
    CONSTRAINT fk_projects_user FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
