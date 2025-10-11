package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/Takanpon2512/english-app/internal/model"
	"github.com/Takanpon2512/english-app/internal/repository"
)

type ProjectService interface {
	CreateProject(userID string, req *model.CreateProjectRequest) (*model.CreateProjectResponse, error)
	GetProjects(userID string, req *model.GetProjectsRequest) (*model.GetProjectsResponse, error)
	GetProjectDetail(userID string, req *model.GetProjectDetailRequest) (*model.GetProjectDetailResponse, error)
}

type projectService struct {
	db   *gorm.DB
	repo repository.ProjectRepository
}

func NewProjectService(db *gorm.DB, repo repository.ProjectRepository) ProjectService {
	return &projectService{
		db:   db,
		repo: repo,
	}
}

func (s *projectService) CreateProject(userID string, req *model.CreateProjectRequest) (*model.CreateProjectResponse, error) {
	now := time.Now()
	project := &model.Project{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {

		if err := s.repo.CreateProject(project); err != nil {
			return fmt.Errorf("プロジェクトの作成に失敗しました: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// レスポンスの作成
	response := &model.CreateProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}

	return response, nil
}

// GetProjects プロジェクト一覧を取得する
func (s *projectService) GetProjects(userID string, req * model.GetProjectsRequest) (*model.GetProjectsResponse, error) {
	return s.repo.GetProjects(req)
}

// GetProjectDetail プロジェクト詳細を取得する
func (s *projectService) GetProjectDetail(userID string, req *model.GetProjectDetailRequest) (*model.GetProjectDetailResponse, error) {
	return s.repo.GetProjectDetail(req)
}
