package repository

import (
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

type SceneRepository struct{ DB *gorm.DB }

func NewSceneRepository(db *gorm.DB) *SceneRepository { return &SceneRepository{DB: db} }

func (r *SceneRepository) Create(s *domain.Scene) error { return r.DB.Create(s).Error }

func (r *SceneRepository) FindByCampaign(campaignID uint) ([]domain.Scene, error) {
	var scenes []domain.Scene
	err := r.DB.Where("campaign_id = ?", campaignID).Order("name").Find(&scenes).Error
	return scenes, err
}

func (r *SceneRepository) FindByID(id uint) (domain.Scene, error) {
	var s domain.Scene
	err := r.DB.Preload("Tokens").First(&s, id).Error
	return s, err
}

func (r *SceneRepository) Update(s *domain.Scene) error { return r.DB.Save(s).Error }

func (r *SceneRepository) Delete(id uint) error { return r.DB.Delete(&domain.Scene{}, id).Error }

func (r *SceneRepository) CreateToken(t *domain.Token) error { return r.DB.Create(t).Error }

func (r *SceneRepository) FindTokensByScene(sceneID uint) ([]domain.Token, error) {
	var tokens []domain.Token
	err := r.DB.Where("scene_id = ?", sceneID).Find(&tokens).Error
	return tokens, err
}

func (r *SceneRepository) FindTokenByID(id uint) (domain.Token, error) {
	var t domain.Token
	err := r.DB.First(&t, id).Error
	return t, err
}

func (r *SceneRepository) UpdateToken(t *domain.Token) error { return r.DB.Save(t).Error }

func (r *SceneRepository) DeleteToken(id uint) error { return r.DB.Delete(&domain.Token{}, id).Error }
