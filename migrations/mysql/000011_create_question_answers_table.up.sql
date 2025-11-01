-- Create question_answers table
CREATE TABLE question_answers (
    id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    project_id CHAR(36) NOT NULL,
    question_template_master_id CHAR(36) NOT NULL,
    user_answer TEXT NULL,
    challenge_count INT NOT NULL DEFAULT 1,
    status VARCHAR(20) NOT NULL DEFAULT 'PROCESSING',
    created_by CHAR(36) NOT NULL,
    updated_by CHAR(36) NOT NULL,
    deleted_by CHAR(36) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);
