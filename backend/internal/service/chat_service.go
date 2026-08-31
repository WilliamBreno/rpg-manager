package service

import (
	"errors"
	"sort"

	"rpg-manager/internal/domain"
)

type chatMessageRepo interface {
	Create(m *domain.ChatMessage) error
	FindByID(id uint) (domain.ChatMessage, error)
	FindByCampaign(campaignID uint, limit int) ([]domain.ChatMessage, error)
}

type ChatService struct{ repo chatMessageRepo }

func NewChatService(repo chatMessageRepo) *ChatService { return &ChatService{repo: repo} }

// Send persiste e recarrega com o remetente preenchido (Preload) — sem isso
// o objeto em memória só tem SenderUserID, e o broadcast via WebSocket (que
// usa esse mesmo objeto, não uma nova consulta) chegaria pro cliente com
// `sender: {}` vazio até a próxima vez que o histórico fosse recarregado via
// REST.
func (s *ChatService) Send(m *domain.ChatMessage) error {
	if m.Text == "" {
		return errors.New("mensagem vazia")
	}
	if err := s.repo.Create(m); err != nil {
		return err
	}
	reloaded, err := s.repo.FindByID(m.ID)
	if err != nil {
		return err
	}
	*m = reloaded
	return nil
}

// History devolve as últimas `limit` mensagens em ordem cronológica
// (mais antiga primeiro) — o repositório busca mais-recentes-primeiro pra
// paginação futura ser natural (LIMIT sem OFFSET pega sempre as últimas),
// então a ordem é invertida aqui só pra exibição.
func (s *ChatService) History(campaignID uint, limit int) ([]domain.ChatMessage, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	messages, err := s.repo.FindByCampaign(campaignID, limit)
	if err != nil {
		return nil, err
	}
	sort.Slice(messages, func(i, j int) bool { return messages[i].ID < messages[j].ID })
	return messages, nil
}
