ALTER TABLE correction_results 
ADD CONSTRAINT fk_correction_results_question_answer_id 
FOREIGN KEY (question_answer_id) REFERENCES question_answers(id) 
ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE correction_results 
ADD CONSTRAINT fk_correction_results_question_template_master_id 
FOREIGN KEY (question_template_master_id) REFERENCES question_template_masters(id) 
ON DELETE RESTRICT ON UPDATE CASCADE;
