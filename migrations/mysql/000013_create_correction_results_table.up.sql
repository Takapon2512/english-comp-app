CREATE TABLE correction_results (
    id CHAR(36) PRIMARY KEY,
    question_answer_id CHAR(36) NOT NULL,
    question_template_master_id CHAR(36) NOT NULL,
    get_points INT NOT NULL DEFAULT 0,
    example_correction TEXT NULL,
    correct_rate INT NULL,
    advice TEXT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PROCESSING',
    challenge_count INT NOT NULL DEFAULT 1,
    created_by CHAR(36) NOT NULL,
    updated_by CHAR(36) NOT NULL,
    deleted_by CHAR(36) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);
