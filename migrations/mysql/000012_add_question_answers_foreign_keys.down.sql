-- Drop foreign key constraints for question_answers table
ALTER TABLE question_answers 
DROP FOREIGN KEY fk_question_answers_user_id;

ALTER TABLE question_answers 
DROP FOREIGN KEY fk_question_answers_project_id;

ALTER TABLE question_answers 
DROP FOREIGN KEY fk_question_answers_question_template_master_id;
