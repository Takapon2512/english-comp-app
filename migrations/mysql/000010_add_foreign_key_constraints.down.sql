-- Remove foreign key constraints for project_questions table
ALTER TABLE project_questions DROP CONSTRAINT fk_project_questions_project_id;
ALTER TABLE project_questions DROP CONSTRAINT fk_project_questions_question_template_master_id;
