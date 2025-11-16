-- 弱点分析テーブル間の外部キー制約を削除

-- weakness_learning_advice テーブルの外部キー制約削除
ALTER TABLE weakness_learning_advice 
DROP FOREIGN KEY fk_weakness_learning_advice_analysis_id;

-- weakness_detailed_analyses テーブルの外部キー制約削除
ALTER TABLE weakness_detailed_analyses 
DROP FOREIGN KEY fk_weakness_detailed_analyses_analysis_id;

-- weakness_category_analyses テーブルの外部キー制約削除
ALTER TABLE weakness_category_analyses 
DROP FOREIGN KEY fk_weakness_category_analyses_category_id;

ALTER TABLE weakness_category_analyses 
DROP FOREIGN KEY fk_weakness_category_analyses_analysis_id;

-- weakness_analyses テーブルの外部キー制約削除
ALTER TABLE weakness_analyses 
DROP FOREIGN KEY fk_weakness_analyses_user_id;

ALTER TABLE weakness_analyses 
DROP FOREIGN KEY fk_weakness_analyses_project_id;
