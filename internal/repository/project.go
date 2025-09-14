package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
)

type ProjectRepository interface {
	CreateProject(project *model.Project) error
	CreateTags(tags []model.Tag) error
	GetProjects(req *model.GetProjectsRequest) (*model.GetProjectsResponse, error)
	GetProjectDetail(req *model.GetProjectDetailRequest) (*model.GetProjectDetailResponse, error)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) CreateProject(project *model.Project) error {
	if project.ID == "" {
		project.ID = uuid.New().String()
	}
	return r.db.Create(project).Error
}

func (r *projectRepository) CreateTags(tags []model.Tag) error {
	if len(tags) == 0 {
		return nil
	}

	for i := range tags {
		if tags[i].ID == "" {
			tags[i].ID = uuid.New().String()
		}
	}

	return r.db.Table("project_tags").Create(&tags).Error
}

// GetProjects プロジェクト一覧を取得する
func (r *projectRepository) GetProjects(req *model.GetProjectsRequest) (*model.GetProjectsResponse, error) {
	var projects []model.Project
	var total int64

	query := r.db.Model(&model.Project{}).Where("user_id = ?", req.UserID)

	if req.Tag != "" {
		query = query.Joins("JOIN project_tags ON projects.id = project_tags.project_id").
			Where("project_tags.name = ?", req.Tag)
	}

	// タグ情報を事前読み込み
	query = query.Preload("Tags", func(db *gorm.DB) *gorm.DB {
		return db.Table("project_tags")
	})

	// 総件数を取得
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("プロジェクト数の取得に失敗しました: %w", err)
	}

	// プロジェクト一覧を取得（作成日時の降順）
	// ページネーションのデバッグ出力
	offset := (req.Page - 1) * req.PerPage

	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(req.PerPage).
		Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("プロジェクト一覧の取得に失敗しました: %w", err)
	}

	return &model.GetProjectsResponse{
		Projects: projects,
		Total:    int(total),
		Page:     req.Page,
		PerPage:  req.PerPage,
	}, nil
}

// GetProjectDetail プロジェクト詳細を取得する
func (r *projectRepository) GetProjectDetail(req *model.GetProjectDetailRequest) (*model.GetProjectDetailResponse, error) {
	var project model.Project

	query := r.db.Model(&model.Project{}).Where("id = ?", req.ID)

	if err := query.Preload("Tags", func(db *gorm.DB) *gorm.DB {
		return db.Table("project_tags")
	}).First(&project).Error; err != nil {
		return nil, fmt.Errorf("プロジェクトの取得に失敗しました: %w", err)
	}

	return &model.GetProjectDetailResponse{
		Project: project,
	}, nil
}
