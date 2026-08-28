package service

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"rpg-manager/internal/domain"
)

type periciaRepo interface {
	GetAll(edition string) ([]domain.Pericia, error)
	GetByCharacter(characterID uint) ([]domain.CharacterPericia, error)
	Save(characterID uint, names []string) error
	SetExpertise(characterID uint, periciaName string, expertise bool) error
}

// PericiaService também recebe DB direto (mesmo padrão de InventoryService/
// BackgroundService) só pra validar o limite de Expertise contra a classe e
// o nível do personagem — não vale a pena criar um CharacterRepo dependency
// só pra isso.
type PericiaService struct {
	repo periciaRepo
	DB   *gorm.DB
}

func NewPericiaService(repo periciaRepo, db *gorm.DB) *PericiaService {
	return &PericiaService{repo: repo, DB: db}
}

func (s *PericiaService) GetAll(edition string) ([]domain.Pericia, error) {
	return s.repo.GetAll(edition)
}
func (s *PericiaService) GetByCharacter(id uint) ([]domain.CharacterPericia, error) {
	return s.repo.GetByCharacter(id)
}
func (s *PericiaService) Save(id uint, names []string) error {
	return s.repo.Save(id, names)
}

// SetExpertise marca/desmarca Especialização numa perícia já proficiente.
// Validação só roda ao MARCAR (desmarcar sempre é permitido, nunca estoura
// limite nenhum): a classe do personagem precisa conceder Expertise no nível
// atual (ExpertiseSlotsFor), e o total de perícias já marcadas não pode
// passar do número de vagas daquele nível.
func (s *PericiaService) SetExpertise(characterID uint, periciaName string, expertise bool) error {
	if expertise {
		var character domain.Character
		if err := s.DB.Preload("Class").First(&character, characterID).Error; err != nil {
			return err
		}
		slots := ExpertiseSlotsFor(character.Class.Name, character.Level)
		if slots == 0 {
			return errors.New("essa classe não concede Especialização neste nível")
		}
		var current int64
		s.DB.Model(&domain.CharacterPericia{}).
			Where("character_id = ? AND expertise = ?", characterID, true).
			Count(&current)
		if int(current) >= slots {
			return fmt.Errorf("limite de %d especialização(ões) já atingido neste nível", slots)
		}
	}
	if err := s.repo.SetExpertise(characterID, periciaName, expertise); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("o personagem precisa ser proficiente nessa perícia antes de ter Especialização nela")
		}
		return err
	}
	return nil
}
