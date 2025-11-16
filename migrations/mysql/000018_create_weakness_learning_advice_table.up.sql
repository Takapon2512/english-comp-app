-- WeaknessLearningAdvice テーブルの作成
-- 個別化された学習アドバイスと提案を管理するテーブル
CREATE TABLE weakness_learning_advice (
    id CHAR(36) PRIMARY KEY COMMENT 'レコードの一意識別子',
    analysis_id CHAR(36) NOT NULL COMMENT '親分析レコードのID（外部キー）',
    
    -- 学習アドバイス情報
    learning_advice TEXT COMMENT '個別学習アドバイス（具体的な学習方法や注意点）',
    recommended_actions JSON COMMENT '推奨アクションのJSON配列（具体的な学習行動の提案）',
    next_goals JSON COMMENT '次の学習目標のJSON配列（短期・中期目標の設定）',
    study_plan TEXT COMMENT '個別学習プラン（期間、内容、方法を含む詳細プラン）',
    motivational_message TEXT COMMENT 'モチベーション向上メッセージ（励ましや成長の認識）',
    
    -- 標準的なデータベース管理フィールド
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'レコード作成日時',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'レコード最終更新日時',
    deleted_at DATETIME NULL COMMENT '論理削除日時',
    created_by CHAR(36) NOT NULL COMMENT 'レコード作成者のユーザーID',
    updated_by CHAR(36) NOT NULL COMMENT 'レコード最終更新者のユーザーID',
    
    -- インデックス
    INDEX idx_weakness_learning_advice_analysis_id (analysis_id),
    INDEX idx_weakness_learning_advice_deleted_at (deleted_at),
    UNIQUE KEY uk_weakness_learning_advice_analysis_id (analysis_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='個別化学習アドバイステーブル';
