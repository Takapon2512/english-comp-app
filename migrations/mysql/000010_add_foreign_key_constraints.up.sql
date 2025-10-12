-- Add foreign key constraints for project_questions table
ALTER TABLE project_questions 
ADD CONSTRAINT fk_project_questions_project_id 
FOREIGN KEY (project_id) REFERENCES projects(id) 
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE project_questions 
ADD CONSTRAINT fk_project_questions_question_template_master_id 
FOREIGN KEY (question_template_master_id) REFERENCES question_template_masters(id) 
ON DELETE CASCADE ON UPDATE CASCADE;
