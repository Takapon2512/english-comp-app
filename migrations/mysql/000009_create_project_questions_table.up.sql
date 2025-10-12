CREATE TABLE project_questions (
    id CHAR(36) PRIMARY KEY,
    project_id CHAR(36) NOT NULL,
    question_template_master_id CHAR(36) NOT NULL,
    created_by CHAR(36) NOT NULL,
    updated_by CHAR(36) NOT NULL,
    deleted_by CHAR(36) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);
