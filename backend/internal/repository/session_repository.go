package repository

import (
	"rpg-manager/internal/domain"

	"gorm.io/gorm"
)

type SessionRepository struct{ DB *gorm.DB }

func NewSessionRepository(db *gorm.DB) *SessionRepository { return &SessionRepository{DB: db} }

func (r *SessionRepository) Create(s *domain.Session) error { return r.DB.Create(s).Error }

func (r *SessionRepository) FindByCampaign(campaignID uint) ([]domain.Session, error) {
	var sessions []domain.Session
	err := r.DB.Where("campaign_id = ?", campaignID).Order("started_at desc").Find(&sessions).Error
	return sessions, err
}

func (r *SessionRepository) FindByID(id uint) (domain.Session, error) {
	var s domain.Session
	err := r.DB.First(&s, id).Error
	return s, err
}

// FindActiveByCampaign retorna a sessão em andamento (EndedAt nulo) de uma
// campanha, se houver — usada pra impedir duas sessões abertas ao mesmo
// tempo e, na Etapa 5, pra saber em qual sala/hub um jogador que conecta cai.
func (r *SessionRepository) FindActiveByCampaign(campaignID uint) (domain.Session, error) {
	var s domain.Session
	err := r.DB.Where("campaign_id = ? AND ended_at IS NULL", campaignID).First(&s).Error
	return s, err
}

func (r *SessionRepository) Update(s *domain.Session) error { return r.DB.Save(s).Error }
