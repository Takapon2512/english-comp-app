CREATE TABLE IF NOT EXISTS project_tags (
    id CHAR(36) NOT NULL,
    project_id CHAR(36) NOT NULL,
    name VARCHAR(30) NOT NULL,
    created_at DATETIME NOT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX project_id_name_idx (project_id, name),
    INDEX name_idx (name),
    CONSTRAINT fk_project_tags_project FOREIGN KEY (project_id)
        REFERENCES projects(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
