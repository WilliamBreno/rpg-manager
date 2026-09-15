package service

import (
	"errors"

	"rpg-manager/internal/domain"
)

type campaignMembershipRepo interface {
	Create(m *domain.CampaignMembership) error
	FindByCampaign(campaignID uint) ([]domain.CampaignMembership, error)
	FindByCampaignAndUser(campaignID, userID uint) (domain.CampaignMembership, error)
	FindByUser(userID uint) ([]domain.CampaignMembership, error)
	FindByID(id uint) (domain.CampaignMembership, error)
	Update(m *domain.CampaignMembership) error
}

type userLookup interface {
	FindByEmail(email string) (domain.User, error)
}

type membershipCharacterLookup interface {
	FindByID(id uint) (domain.Character, error)
}

type CampaignMembershipService struct {
	repo      campaignMembershipRepo
	userRepo  userLookup
	charRepo  membershipCharacterLookup
}

func NewCampaignMembershipService(repo campaignMembershipRepo, userRepo userLookup, charRepo membershipCharacterLookup) *CampaignMembershipService {
	return &CampaignMembershipService{repo: repo, userRepo: userRepo, charRepo: charRepo}
}

// Invite resolve o e-mail pra um usuário existente e cria o convite
// (Status: invited). É sempre o mestre quem convida — não existe fluxo de
// jogador convidando outro jogador pra campanha (ver Etapa 6 do
// SISTEMA_MESTRE.md: "não é o jogador quem convida outros jogadores").
func (s *CampaignMembershipService) Invite(campaignID uint, email string) (domain.CampaignMembership, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return domain.CampaignMembership{}, errors.New("nenhum usuário encontrado com esse e-mail")
	}
	if _, err := s.repo.FindByCampaignAndUser(campaignID, user.ID); err == nil {
		return domain.CampaignMembership{}, errors.New("esse jogador já foi convidado pra essa campanha")
	}
	m := domain.CampaignMembership{CampaignID: campaignID, UserID: user.ID, Status: domain.MembershipInvited}
	if err := s.repo.Create(&m); err != nil {
		return domain.CampaignMembership{}, err
	}
	return m, nil
}

func (s *CampaignMembershipService) GetByCampaign(campaignID uint) ([]domain.CampaignMembership, error) {
	return s.repo.FindByCampaign(campaignID)
}

// GetPendingForUser lista os convites (`invited`) que um jogador ainda não
// respondeu, cruzando todas as campanhas — é a área "adicionar mestre(s)"
// do lado do jogador (Etapa 6): aceitar convite de campanha.
func (s *CampaignMembershipService) GetPendingForUser(userID uint) ([]domain.CampaignMembership, error) {
	all, err := s.repo.FindByUser(userID)
	if err != nil {
		return nil, err
	}
	var pending []domain.CampaignMembership
	for _, m := range all {
		if m.Status == domain.MembershipInvited {
			pending = append(pending, m)
		}
	}
	return pending, nil
}

// GetAcceptedForUser lista as campanhas que o jogador já está de fato
// participando — usado pra ele navegar até a Sala ao vivo de cada uma.
func (s *CampaignMembershipService) GetAcceptedForUser(userID uint) ([]domain.CampaignMembership, error) {
	all, err := s.repo.FindByUser(userID)
	if err != nil {
		return nil, err
	}
	var accepted []domain.CampaignMembership
	for _, m := range all {
		if m.Status == domain.MembershipAccepted {
			accepted = append(accepted, m)
		}
	}
	return accepted, nil
}

func (s *CampaignMembershipService) GetByID(id uint) (domain.CampaignMembership, error) {
	return s.repo.FindByID(id)
}

// Respond é o jogador aceitando ou recusando um convite — só o próprio
// convidado pode responder (checado no handler via UserID). Ao aceitar, é
// obrigatório escolher qual personagem representará o jogador nesta
// campanha (pedido explícito do usuário — antes disso o vínculo de
// personagem era opcional e a maioria dos jogadores entrava sem nenhum
// personagem vinculado, o que quebrava o elenco do mestre e a concessão de
// recompensas/XP). O personagem precisa pertencer ao próprio jogador que
// está respondendo — sem essa checagem, o campo character_id do payload
// permitiria vincular o personagem de outra pessoa à campanha.
func (s *CampaignMembershipService) Respond(m *domain.CampaignMembership, accept bool, characterID *uint) error {
	if m.Status != domain.MembershipInvited {
		return errors.New("esse convite já foi respondido")
	}
	if accept {
		if characterID == nil {
			return errors.New("escolha um personagem para entrar na campanha")
		}
		character, err := s.charRepo.FindByID(*characterID)
		if err != nil || character.UserID != m.UserID {
			return errors.New("personagem inválido")
		}
		m.Status = domain.MembershipAccepted
		m.CharacterID = characterID
	} else {
		m.Status = domain.MembershipDeclined
	}
	return s.repo.Update(m)
}

// IsAcceptedMember diz se userID pode entrar na sala/chat de campaignID —
// usado pelo gate de WebSocket da Etapa 5 (o mestre sempre pode, checado
// separadamente pelo chamador via Campaign.MasterID).
func (s *CampaignMembershipService) IsAcceptedMember(campaignID, userID uint) bool {
	m, err := s.repo.FindByCampaignAndUser(campaignID, userID)
	return err == nil && m.Status == domain.MembershipAccepted
}
