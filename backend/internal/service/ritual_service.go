package service

import (
	"rpg-manager/internal/domain"
	"rpg-manager/internal/repository"

	"gorm.io/gorm"
)

// ritualAccessRule descreve a característica "Conjuração Ritual" de uma
// classe 4e — extraído verbatim do texto de cada classe (ver levantamento
// aprovado pelo usuário antes desta implementação, não de suposição
// genérica sobre "como rituais costumam funcionar").
type ritualAccessRule struct {
	Level1Slots int      // total de rituais conhecidos ao ganhar a característica, no nível 1
	FixedRituals []string // rituais fixos, concedidos automaticamente sem escolha (consomem 1 vaga cada)

	// RequiredPrerequisiteClass: quando não vazio, pelo menos 1 das vagas
	// (além das fixas) deve ser preenchida por um ritual cujo Prerequisite
	// seja essa mesma classe (ex: Bardo — "um ritual de nível 1 que tenha a
	// classe bardo como pré-requisito"). Não é travado no banco, só sinalizado
	// na resposta pra a UI orientar a escolha.
	RequiredPrerequisiteClass string

	// RestrictedOptions: quando não vazio, pelo menos 1 das vagas deve vir
	// desse conjunto específico de nomes (ex: Psionista escolhe entre Disco
	// Flutuante de Tenser e Enviar Mensagem).
	RestrictedOptions []string

	// GrowthLevels: nível do personagem -> vagas adicionais ganhas naquele
	// nível. Só o Mago tem isso confirmado nos livros (5º/11º/15º/21º/25º,
	// +2 cada vez) — as outras 5 classes não têm nenhuma tabela de
	// crescimento do número de rituais conhecidos documentada na
	// característica de classe (podem ganhar mais execuções grátis por dia,
	// mas isso não aumenta quantos rituais estão no livro de rituais).
	GrowthLevels map[int]int
}

// ritualAccessRules4e — as 6 classes confirmadas no levantamento (LJ1: Mago,
// Clérigo; LJ2: Bardo, Druida, Invocador; LJ3: Psionista). Nenhuma outra
// classe 4e tem Conjuração Ritual como característica nativa, e nenhuma
// raça concede a característica ou um ritual específico (achado do
// levantamento, único ponto relacionado é um talento de paragon exclusivo
// de githzerai que nem sequer concede a característica em si — fora do
// escopo desta modelagem de acesso por classe).
var ritualAccessRules4e = map[string]ritualAccessRule{
	"Mago": {
		Level1Slots:  3,
		GrowthLevels: map[int]int{5: 2, 11: 2, 15: 2, 21: 2, 25: 2},
	},
	"Clérigo": {
		Level1Slots:  2,
		FixedRituals: []string{"Repouso Tranquilo"},
	},
	"Bardo": {
		Level1Slots:               2,
		RequiredPrerequisiteClass: "Bardo",
	},
	"Druida": {
		Level1Slots:  2,
		FixedRituals: []string{"Mensageiro Animal"},
	},
	"Invocador": {
		Level1Slots:  2,
		FixedRituals: []string{"Mão do Destino"},
	},
	"Psionista": {
		Level1Slots:       2,
		RestrictedOptions: []string{"Disco Flutuante de Tenser", "Enviar Mensagem"},
	},
}

// totalRitualSlots soma o nível 1 + qualquer crescimento em níveis <= level.
func totalRitualSlots(rule ritualAccessRule, level int) int {
	total := rule.Level1Slots
	for atLevel, bonus := range rule.GrowthLevels {
		if level >= atLevel {
			total += bonus
		}
	}
	return total
}

type RitualAccessInfo struct {
	HasRitualCasting          bool            `json:"has_ritual_casting"`
	ClassName                 string          `json:"class_name"`
	Known                     []domain.Ritual `json:"known"`
	TotalSlots                int             `json:"total_slots"`
	RemainingChoices          int             `json:"remaining_choices"`
	RequiredPrerequisiteClass string          `json:"required_prerequisite_class,omitempty"`
	RestrictedOptions         []string        `json:"restricted_options,omitempty"`
	AvailableRituals          []domain.Ritual `json:"available_rituals"`
}

type ritualCharacterRepo interface {
	FindByID(id uint) (domain.Character, error)
}

type RitualService struct {
	repo    *repository.RitualRepository
	charRepo ritualCharacterRepo
	db      *gorm.DB
}

func NewRitualService(repo *repository.RitualRepository, charRepo ritualCharacterRepo, db *gorm.DB) *RitualService {
	return &RitualService{repo: repo, charRepo: charRepo, db: db}
}

func (s *RitualService) GetAll(edition string) ([]domain.Ritual, error) {
	return s.repo.GetAll(edition)
}
func (s *RitualService) GetByCharacter(id uint) ([]domain.Ritual, error) {
	return s.repo.GetByCharacter(id)
}
func (s *RitualService) Add(charID, ritualID uint) error {
	return s.repo.Add(charID, ritualID)
}
func (s *RitualService) Remove(charID, ritualID uint) error {
	return s.repo.Remove(charID, ritualID)
}

// GetAccess calcula, pra um personagem específico, se a classe dele tem
// Conjuração Ritual e, se tiver, quais rituais fixos ele já deveria ter
// (concedendo os que faltarem — cobre tanto personagens novos quanto
// personagens 4e criados antes desta feature existir), quantas vagas ainda
// pode escolher, e a lista de rituais disponíveis pra escolha (nível do
// ritual <= nível do personagem, e sem Pré-requisito de outra classe).
func (s *RitualService) GetAccess(characterID uint) (RitualAccessInfo, error) {
	character, err := s.charRepo.FindByID(characterID)
	if err != nil {
		return RitualAccessInfo{}, err
	}

	rule, ok := ritualAccessRules4e[character.Class.Name]
	if !ok || character.Edition != "4e" {
		return RitualAccessInfo{HasRitualCasting: false, ClassName: character.Class.Name}, nil
	}

	// Garante que os rituais fixos da classe já estejam concedidos —
	// idempotente, cobre personagens criados antes desta feature existir.
	knownNames := map[string]bool{}
	for _, r := range character.Rituals {
		knownNames[r.Name] = true
	}
	for _, fixedName := range rule.FixedRituals {
		if knownNames[fixedName] {
			continue
		}
		fixed, err := s.repo.FindByName(fixedName, "4e")
		if err == nil {
			s.repo.Add(character.ID, fixed.ID)
			knownNames[fixedName] = true
		}
	}

	known, err := s.repo.GetByCharacter(character.ID)
	if err != nil {
		return RitualAccessInfo{}, err
	}

	totalSlots := totalRitualSlots(rule, character.Level)
	remaining := totalSlots - len(known)
	if remaining < 0 {
		remaining = 0
	}

	allRituals, err := s.repo.GetAll("4e")
	if err != nil {
		return RitualAccessInfo{}, err
	}
	available := make([]domain.Ritual, 0, len(allRituals))
	for _, r := range allRituals {
		if knownNames[r.Name] {
			continue
		}
		if r.Level > character.Level {
			continue
		}
		if r.Prerequisite != "" && r.Prerequisite != character.Class.Name {
			continue
		}
		available = append(available, r)
	}

	return RitualAccessInfo{
		HasRitualCasting:          true,
		ClassName:                 character.Class.Name,
		Known:                     known,
		TotalSlots:                totalSlots,
		RemainingChoices:          remaining,
		RequiredPrerequisiteClass: rule.RequiredPrerequisiteClass,
		RestrictedOptions:         rule.RestrictedOptions,
		AvailableRituals:          available,
	}, nil
}
