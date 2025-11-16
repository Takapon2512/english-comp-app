-- 弱点分析テーブル間の外部キー制約を追加

-- weakness_analyses テーブルの外部キー制約
ALTER TABLE weakness_analyses 
ADD CONSTRAINT fk_weakness_analyses_project_id 
FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;

ALTER TABLE weakness_analyses 
ADD CONSTRAINT fk_weakness_analyses_user_id 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- weakness_category_analyses テーブルの外部キー制約
ALTER TABLE weakness_category_analyses 
ADD CONSTRAINT fk_weakness_category_analyses_analysis_id 
FOREIGN KEY (analysis_id) REFERENCES weakness_analyses(id) ON DELETE CASCADE;

ALTER TABLE weakness_category_analyses 
ADD CONSTRAINT fk_weakness_category_analyses_category_id 
FOREIGN KEY (category_id) REFERENCES category_masters(id) ON DELETE CASCADE;

-- weakness_detailed_analyses テーブルの外部キー制約
ALTER TABLE weakness_detailed_analyses 
ADD CONSTRAINT fk_weakness_detailed_analyses_analysis_id 
FOREIGN KEY (analysis_id) REFERENCES weakness_analyses(id) ON DELETE CASCADE;

-- weakness_learning_advice テーブルの外部キー制約
ALTER TABLE weakness_learning_advice 
ADD CONSTRAINT fk_weakness_learning_advice_analysis_id 
FOREIGN KEY (analysis_id) REFERENCES weakness_analyses(id) ON DELETE CASCADE;
