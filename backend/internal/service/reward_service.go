package service

import (
	"errors"

	"rpg-manager/internal/domain"
)

type rewardRepo interface {
	Create(reward *domain.Reward) error
	FindByCampaign(campaignID uint) ([]domain.Reward, error)
	FindByCharacter(characterID uint) ([]domain.Reward, error)
}

type magicItemRepo interface {
	Create(m *domain.MagicItem) error
	FindByCampaign(campaignID uint) ([]domain.MagicItem, error)
	FindByID(id uint) (domain.MagicItem, error)
	Delete(id uint) error
}

// RewardService concede moeda/item a UM personagem por chamada — "dar a
// todos" na UI (Etapa 8) é o handler chamando isso uma vez por destinatário,
// gerando uma linha de Reward por personagem (ver comentário em
// domain.Reward: 1 linha = 1 entrega a 1 personagem, nunca N destinatários
// numa linha só).
type RewardService struct {
	repo      rewardRepo
	itemRepo  magicItemRepo
	inventory *InventoryService
}

func NewRewardService(repo rewardRepo, itemRepo magicItemRepo, inventory *InventoryService) *RewardService {
	return &RewardService{repo: repo, itemRepo: itemRepo, inventory: inventory}
}

// GrantCurrency credita as moedas de verdade no personagem (via
// InventoryService.GrantCurrency, soma ao que já tem) e grava o histórico.
func (s *RewardService) GrantCurrency(campaignID, characterID, grantedBy uint, cp, sp, ep, gp, pp int, note string) (domain.Reward, error) {
	if cp < 0 || sp < 0 || ep < 0 || gp < 0 || pp < 0 {
		return domain.Reward{}, errors.New("valores de moeda não podem ser negativos")
	}
	if cp == 0 && sp == 0 && ep == 0 && gp == 0 && pp == 0 {
		return domain.Reward{}, errors.New("informe ao menos um valor de moeda maior que zero")
	}
	if _, err := s.inventory.GrantCurrency(characterID, cp, sp, ep, gp, pp); err != nil {
		return domain.Reward{}, err
	}
	reward := domain.Reward{
		CampaignID: campaignID, CharacterID: characterID, GrantedByUserID: grantedBy, Kind: domain.RewardCurrency,
		CopperPieces: cp, SilverPieces: sp, ElectrumPieces: ep, GoldPieces: gp, PlatinumPieces: pp, Note: note,
	}
	if err := s.repo.Create(&reward); err != nil {
		return domain.Reward{}, err
	}
	return reward, nil
}

// GrantItem só registra a entrega (ver domain.MagicItem — não há catálogo de
// inventário estruturado pra item de campanha, os próprios registros de
// Reward tipo "item" de um personagem já servem como sua lista de itens de
// campanha recebidos, sem precisar de uma tabela de posse separada).
func (s *RewardService) GrantItem(campaignID, characterID, grantedBy, magicItemID uint, note string) (domain.Reward, error) {
	item, err := s.itemRepo.FindByID(magicItemID)
	if err != nil || item.CampaignID != campaignID {
		return domain.Reward{}, errors.New("item mágico não pertence a essa campanha")
	}
	reward := domain.Reward{
		CampaignID: campaignID, CharacterID: characterID, GrantedByUserID: grantedBy,
		Kind: domain.RewardItem, MagicItemID: &magicItemID, Note: note,
	}
	if err := s.repo.Create(&reward); err != nil {
		return domain.Reward{}, err
	}
	return reward, nil
}

func (s *RewardService) GetByCampaign(campaignID uint) ([]domain.Reward, error) {
	return s.repo.FindByCampaign(campaignID)
}

func (s *RewardService) GetByCharacter(characterID uint) ([]domain.Reward, error) {
	return s.repo.FindByCharacter(characterID)
}

// MagicItemService — catálogo de itens de campanha, CRUD simples.
type MagicItemService struct{ repo magicItemRepo }

func NewMagicItemService(repo magicItemRepo) *MagicItemService { return &MagicItemService{repo: repo} }

func (s *MagicItemService) Create(m *domain.MagicItem) error {
	if m.Name == "" {
		return errors.New("nome do item é obrigatório")
	}
	return s.repo.Create(m)
}

func (s *MagicItemService) GetByCampaign(campaignID uint) ([]domain.MagicItem, error) {
	return s.repo.FindByCampaign(campaignID)
}

func (s *MagicItemService) GetByID(id uint) (domain.MagicItem, error) { return s.repo.FindByID(id) }

func (s *MagicItemService) Delete(id uint) error { return s.repo.Delete(id) }
