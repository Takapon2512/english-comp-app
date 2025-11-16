-- WeaknessDetailedAnalysis テーブルの作成
-- 4つの主要学習領域の詳細分析結果を管理するテーブル
CREATE TABLE weakness_detailed_analyses (
    id CHAR(36) PRIMARY KEY COMMENT 'レコードの一意識別子',
    analysis_id CHAR(36) NOT NULL COMMENT '親分析レコードのID（外部キー）',
    
    -- 文法領域の分析結果
    grammar_score INT NOT NULL COMMENT '文法スコア（0-100）',
    grammar_description TEXT COMMENT '文法面の詳細分析説明',
    grammar_examples JSON COMMENT '文法の具体例JSON配列',
    
    -- 語彙領域の分析結果
    vocabulary_score INT NOT NULL COMMENT '語彙スコア（0-100）',
    vocabulary_description TEXT COMMENT '語彙面の詳細分析説明',
    vocabulary_examples JSON COMMENT '語彙の具体例JSON配列',
    
    -- 表現領域の分析結果
    expression_score INT NOT NULL COMMENT '表現スコア（0-100）',
    expression_description TEXT COMMENT '表現面の詳細分析説明',
    expression_examples JSON COMMENT '表現の具体例JSON配列',
    
    -- 構成領域の分析結果
    structure_score INT NOT NULL COMMENT '構成スコア（0-100）',
    structure_description TEXT COMMENT '構成面の詳細分析説明',
    structure_examples JSON COMMENT '構成の具体例JSON配列',
    
    -- 標準的なデータベース管理フィールド
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'レコード作成日時',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'レコード最終更新日時',
    deleted_at DATETIME NULL COMMENT '論理削除日時',
    created_by CHAR(36) NOT NULL COMMENT 'レコード作成者のユーザーID',
    updated_by CHAR(36) NOT NULL COMMENT 'レコード最終更新者のユーザーID',
    
    -- インデックス
    INDEX idx_weakness_detailed_analyses_analysis_id (analysis_id),
    INDEX idx_weakness_detailed_analyses_deleted_at (deleted_at),
    UNIQUE KEY uk_weakness_detailed_analyses_analysis_id (analysis_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='4領域詳細分析結果テーブル';
