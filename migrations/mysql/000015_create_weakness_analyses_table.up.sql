-- WeaknessAnalysis テーブルの作成
-- プロジェクト内でのユーザーの学習弱点分析結果のメインエンティティ
CREATE TABLE weakness_analyses (
    id CHAR(36) PRIMARY KEY COMMENT '分析結果の一意識別子（UUID）',
    project_id CHAR(36) NOT NULL COMMENT '分析対象プロジェクトのID',
    user_id CHAR(36) NOT NULL COMMENT '分析対象ユーザーのID',
    
    -- LLM分析結果の基本情報
    analysis_status VARCHAR(20) NOT NULL DEFAULT 'PROCESSING' COMMENT '分析処理状況（PROCESSING: 処理中, COMPLETED: 完了, FAILED: 失敗）',
    overall_score INT DEFAULT 0 COMMENT 'LLMが算出した総合学習スコア（0-100点）',
    improvement_rate INT DEFAULT 0 COMMENT '前回分析からの改善率（パーセンテージ、-100〜+100）',
    
    -- 分析処理に関するメタデータ
    analysis_date DATETIME NOT NULL COMMENT '分析実行日時',
    analyzed_answers INT NOT NULL DEFAULT 0 COMMENT '分析対象となった回答数',
    data_period_start DATETIME NOT NULL COMMENT '分析対象データの期間開始日',
    data_period_end DATETIME NOT NULL COMMENT '分析対象データの期間終了日',
    llm_model VARCHAR(50) NOT NULL COMMENT '使用したLLMモデル名（例: claude-3-sonnet-20240229）',
    analysis_version VARCHAR(10) NOT NULL COMMENT '分析ロジックのバージョン（プロンプトや処理の変更管理用）',
    
    -- 標準的なデータベース管理フィールド
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'レコード作成日時',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'レコード最終更新日時',
    deleted_at DATETIME NULL COMMENT '論理削除日時（ソフトデリート用）',
    deleted_by CHAR(36) NULL COMMENT '削除実行者のユーザーID',
    created_by CHAR(36) NOT NULL COMMENT 'レコード作成者のユーザーID',
    updated_by CHAR(36) NOT NULL COMMENT 'レコード最終更新者のユーザーID',
    
    -- インデックス
    INDEX idx_weakness_analyses_project_id (project_id),
    INDEX idx_weakness_analyses_user_id (user_id),
    INDEX idx_weakness_analyses_analysis_date (analysis_date),
    INDEX idx_weakness_analyses_deleted_at (deleted_at),
    INDEX idx_weakness_analyses_status (analysis_status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学習弱点分析結果のメインテーブル';
