package service

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"rpg-manager/internal/domain"
)

type sessionRepo interface {
	Create(s *domain.Session) error
	FindByCampaign(campaignID uint) ([]domain.Session, error)
	FindByID(id uint) (domain.Session, error)
	FindActiveByCampaign(campaignID uint) (domain.Session, error)
	Update(s *domain.Session) error
}

type SessionService struct{ repo sessionRepo }

func NewSessionService(repo sessionRepo) *SessionService { return &SessionService{repo: repo} }

// Start abre uma nova sessão pra campanha. Recusa se já existir uma sessão
// em andamento (EndedAt nulo) — evita duas sessões "abertas" simultâneas
// pra mesma campanha, o que confundiria a Etapa 5 (em qual sala um jogador
// que conecta deveria cair).
func (s *SessionService) Start(campaignID uint) (domain.Session, error) {
	_, err := s.repo.FindActiveByCampaign(campaignID)
	if err == nil {
		return domain.Session{}, errors.New("já existe uma sessão em andamento pra essa campanha")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Session{}, err
	}
	session := domain.Session{CampaignID: campaignID, StartedAt: time.Now()}
	if err := s.repo.Create(&session); err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

func (s *SessionService) GetByCampaign(campaignID uint) ([]domain.Session, error) {
	return s.repo.FindByCampaign(campaignID)
}

func (s *SessionService) GetByID(id uint) (domain.Session, error) { return s.repo.FindByID(id) }

// End encerra a sessão e opcionalmente atualiza o resumo/diário no mesmo
// passo (o mestre normalmente escreve o resumo final ao encerrar).
func (s *SessionService) End(session *domain.Session, summary string) error {
	if session.EndedAt != nil {
		return errors.New("sessão já está encerrada")
	}
	now := time.Now()
	session.EndedAt = &now
	if summary != "" {
		session.Summary = summary
	}
	return s.repo.Update(session)
}

// UpdateSummary atualiza só o diário — o mestre pode registrar o que
// aconteceu durante a sessão, não só ao encerrar.
func (s *SessionService) UpdateSummary(session *domain.Session, summary string) error {
	session.Summary = summary
	return s.repo.Update(session)
}

// SetMusic troca a música de fundo da sessão (Etapa 9) — o broadcast do
// evento pros jogadores conectados fica no handler (que tem o *ws.Manager),
// aqui só persiste.
func (s *SessionService) SetMusic(session *domain.Session, musicURL string, playing bool) error {
	if playing {
		session.MusicURL = musicURL
	} else {
		session.MusicURL = ""
	}
	return s.repo.Update(session)
}

// SetActiveScene troca o cenário que a sessão aponta como ativo — é o que a
// Etapa 5 (WebSocket) vai transmitir aos jogadores conectados como
// `scene_changed`. A validação de que sceneID pertence à mesma campanha da
// sessão fica no handler (que já tem acesso ao SceneService).
func (s *SessionService) SetActiveScene(session *domain.Session, sceneID uint) error {
	session.ActiveSceneID = &sceneID
	return s.repo.Update(session)
}
