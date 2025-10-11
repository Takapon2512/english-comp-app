CREATE TABLE question_template_masters (
    id CHAR(36) PRIMARY KEY,
    category_id CHAR(36) NOT NULL,
    question_type VARCHAR(10) NOT NULL,
    english TEXT NOT NULL,
    japanese TEXT NOT NULL,
    status VARCHAR(10) NOT NULL,
    level VARCHAR(10) NOT NULL,
    estimated_time INT NOT NULL,
    points INT NOT NULL,
    created_by CHAR(36) NOT NULL,
    updated_by CHAR(36) NOT NULL,
    deleted_by CHAR(36) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);
