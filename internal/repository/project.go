package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
)

type ProjectRepository interface {
	CreateProject(project *model.Project) error
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

// GetProjects プロジェクト一覧を取得する
func (r *projectRepository) GetProjects(req *model.GetProjectsRequest) (*model.GetProjectsResponse, error) {
	var projects []model.Project
	var total int64

	query := r.db.Model(&model.Project{}).Where("user_id = ?", req.UserID)

	if req.Tag != "" {
		query = query.Joins("JOIN project_tags ON projects.id = project_tags.project_id").
			Joins("JOIN user_tags ON project_tags.user_tags_id = user_tags.id").
			Where("user_tags.name = ?", req.Tag)
	}
	// プロジェクトに紐づく質問数を取得
	query = query.Joins("LEFT JOIN project_questions ON projects.id = project_questions.project_id").
		Select("projects.*, COUNT(project_questions.id) as total_questions").
		Group("projects.id").
		Where("project_questions.deleted_at IS NULL")

	// タグ情報を事前読み込み
	query = query.Preload("Tags").Preload("Tags.UserTag")

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

	// プロジェクトに紐づく質問数を取得
	for i, project := range projects {
		projects[i].TotalQuestions = int(project.TotalQuestions)
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

	query := r.db.Model(&model.Project{}).Where("projects.id = ?", req.ID)

	// プロジェクトに紐づく質問数を取得
	query = query.Joins("LEFT JOIN project_questions ON projects.id = project_questions.project_id").
		Select("projects.*, COUNT(project_questions.id) as total_questions").
		Group("projects.id").
		Where("project_questions.deleted_at IS NULL")

	if err := query.Preload("Tags").Preload("Tags.UserTag").First(&project).Error; err != nil {
		return nil, fmt.Errorf("プロジェクトの取得に失敗しました: %w", err)
	}

	return &model.GetProjectDetailResponse{
		Project: project,
	}, nil
}
