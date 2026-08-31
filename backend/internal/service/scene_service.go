package service

import (
	"errors"

	"rpg-manager/internal/domain"
)

type sceneRepo interface {
	Create(s *domain.Scene) error
	FindByCampaign(campaignID uint) ([]domain.Scene, error)
	FindByID(id uint) (domain.Scene, error)
	Update(s *domain.Scene) error
	Delete(id uint) error

	CreateToken(t *domain.Token) error
	FindTokensByScene(sceneID uint) ([]domain.Token, error)
	FindTokenByID(id uint) (domain.Token, error)
	UpdateToken(t *domain.Token) error
	DeleteToken(id uint) error
}

type SceneService struct{ repo sceneRepo }

func NewSceneService(repo sceneRepo) *SceneService { return &SceneService{repo: repo} }

func (s *SceneService) Create(scene *domain.Scene) error {
	if scene.Name == "" {
		return errors.New("nome do cenário é obrigatório")
	}
	return s.repo.Create(scene)
}

func (s *SceneService) GetByCampaign(campaignID uint) ([]domain.Scene, error) {
	return s.repo.FindByCampaign(campaignID)
}

func (s *SceneService) GetByID(id uint) (domain.Scene, error) { return s.repo.FindByID(id) }

func (s *SceneService) Update(scene *domain.Scene) error {
	if scene.Name == "" {
		return errors.New("nome do cenário é obrigatório")
	}
	return s.repo.Update(scene)
}

func (s *SceneService) Delete(id uint) error { return s.repo.Delete(id) }

// CreateToken adiciona um token decorativo/manual — só posição (X/Y) + label
// + imagem opcional. Sem lógica de movimento por turno (limite de escopo já
// combinado, ver domain.Token).
func (s *SceneService) CreateToken(t *domain.Token) error {
	if t.Label == "" {
		return errors.New("label do token é obrigatório")
	}
	return s.repo.CreateToken(t)
}

func (s *SceneService) GetTokenByID(id uint) (domain.Token, error) { return s.repo.FindTokenByID(id) }

// MoveToken só atualiza X/Y — é o que o arrastar-e-soltar do Konva chama.
func (s *SceneService) MoveToken(t *domain.Token, x, y float64) error {
	t.X, t.Y = x, y
	return s.repo.UpdateToken(t)
}

func (s *SceneService) DeleteToken(id uint) error { return s.repo.DeleteToken(id) }
