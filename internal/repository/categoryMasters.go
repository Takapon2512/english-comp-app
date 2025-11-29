package repository

import (
	"fmt"

	"github.com/Takanpon2512/english-app/internal/model"
	"gorm.io/gorm"
)

type CategoryMastersRepository interface {
	GetCategoryMasters(req *model.GetCategoryMastersSearchRequest) (*model.GetCategoryMastersSearchResponse, error)
	GetCategoryMastersByID(id string) (*model.GetCategoryMastersByIDResponse, error)
	GetCategoryMastersByName(name string) (*model.GetCategoryMastersByIDResponse, error)
}

type categoryMastersRepository struct {
	db *gorm.DB
}

func NewCategoryMastersRepository(db *gorm.DB) CategoryMastersRepository {
	return &categoryMastersRepository{db: db}
}

// GetCategoryMasters カテゴリマスター一覧を取得する
func (r *categoryMastersRepository) GetCategoryMasters(req *model.GetCategoryMastersSearchRequest) (*model.GetCategoryMastersSearchResponse, error) {
	var categoryMasters []model.CategoryMasters
	var total int64

	query := r.db.Model(&model.CategoryMasters{})

	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("カテゴリマスターの取得に失敗しました: %w", err)
	}

	offset := (req.Page - 1) * req.PerPage

	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(req.PerPage).
		Find(&categoryMasters).Error; err != nil {
		return nil, fmt.Errorf("カテゴリマスターの取得に失敗しました: %w", err)
	}

	categoryMastersSummary := make([]model.CategoryMastersSummary, len(categoryMasters))

	for i, categoryMaster := range categoryMasters {
		categoryMastersSummary[i] = model.CategoryMastersSummary{
			ID:   categoryMaster.ID,
			Name: categoryMaster.Name,
		}
	}

	return &model.GetCategoryMastersSearchResponse{
		CategoryMasters: categoryMastersSummary,
		Total:           int(total),
		Page:            req.Page,
		PerPage:         req.PerPage,
	}, nil
}

// GetCategoryMastersByID カテゴリマスターをIDで取得する
func (r *categoryMastersRepository) GetCategoryMastersByID(id string) (*model.GetCategoryMastersByIDResponse, error) {
	var categoryMaster model.CategoryMastersSummary

	if err := r.db.Model(&model.CategoryMasters{}).Where("id = ?", id).First(&categoryMaster).Error; err != nil {
		return nil, fmt.Errorf("カテゴリマスターの取得に失敗しました: %w", err)
	}

	return &model.GetCategoryMastersByIDResponse{
		CategoryMasters: categoryMaster,
	}, nil
}

// GetCategoryMastersByName カテゴリマスターを名前で取得する
func (r *categoryMastersRepository) GetCategoryMastersByName(name string) (*model.GetCategoryMastersByIDResponse, error) {
	var categoryMaster model.CategoryMastersSummary

	if err := r.db.Model(&model.CategoryMasters{}).Where("name = ?", name).First(&categoryMaster).Error; err != nil {
		return nil, fmt.Errorf("カテゴリマスターの取得に失敗しました: %w", err)
	}

	return &model.GetCategoryMastersByIDResponse{
		CategoryMasters: categoryMaster,
	}, nil
}
