package service

import (
	"errors"

	"rpg-manager/internal/domain"
)

type enemyRepo interface {
	Create(e *domain.Enemy) error
	FindByCampaign(campaignID uint) ([]domain.Enemy, error)
	FindByID(id uint) (domain.Enemy, error)
	Update(e *domain.Enemy) error
	Delete(id uint) error

	CreateAbility(a *domain.EnemyAbility) error
	FindAbilityByID(id uint) (domain.EnemyAbility, error)
	UpdateAbility(a *domain.EnemyAbility) error
	DeleteAbility(id uint) error

	CreateLine(l *domain.EnemyLine) error
	FindLineByID(id uint) (domain.EnemyLine, error)
	UpdateLine(l *domain.EnemyLine) error
	DeleteLine(id uint) error
}

type EnemyService struct{ repo enemyRepo }

func NewEnemyService(repo enemyRepo) *EnemyService { return &EnemyService{repo: repo} }

var validEnemyKinds = map[domain.EnemyKind]bool{
	domain.EnemyKindEnemy:   true,
	domain.EnemyKindBoss:    true,
	domain.EnemyKindVillain: true,
}

// Create valida o Enemy/Boss/Vilão e cada uma de suas Abilities (notação de
// dado real, ver dice_notation.go) antes de persistir. Retorna avisos
// (não-bloqueantes) de dano fora da faixa sugerida pro ND informado, um por
// habilidade que gerou aviso — a criação segue normalmente mesmo com avisos,
// só o mestre decide se ajusta ou não (mesmo espírito do próprio Guia do
// Mestre: faixa é sugestão, não regra travada).
func (s *EnemyService) Create(e *domain.Enemy) ([]string, error) {
	if e.Name == "" {
		return nil, errors.New("nome é obrigatório")
	}
	if e.Kind == "" {
		e.Kind = domain.EnemyKindEnemy
	}
	if !validEnemyKinds[e.Kind] {
		return nil, errors.New("kind inválido — use enemy, boss ou villain")
	}
	warnings, err := validateAbilities(e.Abilities, e.ChallengeRating)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(e); err != nil {
		return nil, err
	}
	return warnings, nil
}

func validateAbilities(abilities []domain.EnemyAbility, cr string) ([]string, error) {
	var warnings []string
	for _, a := range abilities {
		if a.Name == "" {
			return nil, errors.New("toda habilidade precisa de um nome")
		}
		_, warning, err := ValidateAbilityDamage(a.Damage, cr)
		if err != nil {
			return nil, err
		}
		if warning != "" {
			warnings = append(warnings, a.Name+": "+warning)
		}
	}
	return warnings, nil
}

func (s *EnemyService) GetByCampaign(campaignID uint) ([]domain.Enemy, error) {
	return s.repo.FindByCampaign(campaignID)
}

func (s *EnemyService) GetByID(id uint) (domain.Enemy, error) { return s.repo.FindByID(id) }

func (s *EnemyService) Update(e *domain.Enemy) error {
	if e.Name == "" {
		return errors.New("nome é obrigatório")
	}
	if !validEnemyKinds[e.Kind] {
		return errors.New("kind inválido — use enemy, boss ou villain")
	}
	return s.repo.Update(e)
}

func (s *EnemyService) Delete(id uint) error { return s.repo.Delete(id) }

// CreateAbility valida notação de dado real (contra o ND do Enemy dono, se
// tiver um) antes de persistir uma habilidade adicionada depois da criação.
func (s *EnemyService) CreateAbility(a *domain.EnemyAbility) (string, error) {
	if a.Name == "" {
		return "", errors.New("nome da habilidade é obrigatório")
	}
	enemy, err := s.repo.FindByID(a.EnemyID)
	if err != nil {
		return "", errors.New("inimigo não encontrado")
	}
	_, warning, err := ValidateAbilityDamage(a.Damage, enemy.ChallengeRating)
	if err != nil {
		return "", err
	}
	if err := s.repo.CreateAbility(a); err != nil {
		return "", err
	}
	return warning, nil
}

func (s *EnemyService) GetAbilityByID(id uint) (domain.EnemyAbility, error) {
	return s.repo.FindAbilityByID(id)
}

func (s *EnemyService) UpdateAbility(a *domain.EnemyAbility) (string, error) {
	if a.Name == "" {
		return "", errors.New("nome da habilidade é obrigatório")
	}
	enemy, err := s.repo.FindByID(a.EnemyID)
	if err != nil {
		return "", errors.New("inimigo não encontrado")
	}
	_, warning, err := ValidateAbilityDamage(a.Damage, enemy.ChallengeRating)
	if err != nil {
		return "", err
	}
	if err := s.repo.UpdateAbility(a); err != nil {
		return "", err
	}
	return warning, nil
}

func (s *EnemyService) DeleteAbility(id uint) error { return s.repo.DeleteAbility(id) }

func (s *EnemyService) CreateLine(l *domain.EnemyLine) error {
	if l.Text == "" {
		return errors.New("texto da fala é obrigatório")
	}
	return s.repo.CreateLine(l)
}

func (s *EnemyService) GetLineByID(id uint) (domain.EnemyLine, error) {
	return s.repo.FindLineByID(id)
}

func (s *EnemyService) UpdateLine(l *domain.EnemyLine) error { return s.repo.UpdateLine(l) }

func (s *EnemyService) DeleteLine(id uint) error { return s.repo.DeleteLine(id) }
