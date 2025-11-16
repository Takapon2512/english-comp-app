-- WeaknessCategoryAnalysis テーブルの作成
-- カテゴリ別の分析結果を管理するテーブル
CREATE TABLE weakness_category_analyses (
    id CHAR(36) PRIMARY KEY COMMENT 'レコードの一意識別子',
    analysis_id CHAR(36) NOT NULL COMMENT '親分析レコードのID（外部キー）',
    category_id CHAR(36) NOT NULL COMMENT 'カテゴリマスターのID',
    
    -- カテゴリ分析結果
    category_name VARCHAR(100) NOT NULL COMMENT 'カテゴリ名（キャッシュ用）',
    score INT NOT NULL COMMENT 'このカテゴリでのスコア（0-100）',
    is_weakness BOOLEAN NOT NULL COMMENT '苦手カテゴリかどうか',
    is_strength BOOLEAN NOT NULL COMMENT '得意カテゴリかどうか',
    
    -- 分析詳細（JSON形式）
    issues JSON COMMENT '具体的な問題点のJSON配列',
    strengths JSON COMMENT '具体的な強みのJSON配列',
    examples JSON COMMENT '具体例のJSON配列',
    
    -- 標準的なデータベース管理フィールド
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'レコード作成日時',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'レコード最終更新日時',
    deleted_at DATETIME NULL COMMENT '論理削除日時',
    created_by CHAR(36) NOT NULL COMMENT 'レコード作成者のユーザーID',
    updated_by CHAR(36) NOT NULL COMMENT 'レコード最終更新者のユーザーID',
    
    -- インデックス
    INDEX idx_weakness_category_analyses_analysis_id (analysis_id),
    INDEX idx_weakness_category_analyses_category_id (category_id),
    INDEX idx_weakness_category_analyses_deleted_at (deleted_at),
    INDEX idx_weakness_category_analyses_weakness (is_weakness),
    INDEX idx_weakness_category_analyses_strength (is_strength)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='カテゴリ別分析結果テーブル';
