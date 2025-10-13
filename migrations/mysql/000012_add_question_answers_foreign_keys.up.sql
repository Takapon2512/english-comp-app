-- Add foreign key constraints for question_answers table
ALTER TABLE question_answers 
ADD CONSTRAINT fk_question_answers_user_id 
FOREIGN KEY (user_id) REFERENCES users(id) 
ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE question_answers 
ADD CONSTRAINT fk_question_answers_project_id 
FOREIGN KEY (project_id) REFERENCES projects(id) 
ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE question_answers 
ADD CONSTRAINT fk_question_answers_question_template_master_id 
FOREIGN KEY (question_template_master_id) REFERENCES question_template_masters(id) 
ON DELETE RESTRICT ON UPDATE CASCADE;
