package domain

import "gorm.io/gorm"

type Antecedent struct {
    gorm.Model
    Name               string `json:"name"`
    Edition            string `json:"edition"`
    Description        string `json:"description"`
    SkillProficiencies string `json:"skill_proficiencies"`
    ToolProficiencies  string `json:"tool_proficiencies"`
    Languages          string `json:"languages"`
    Equipment          string `json:"equipment"`
    Feature            string `json:"feature"`
    FeatureDescription string `json:"feature_description"`
    IsDefault          bool   `json:"is_default" gorm:"default:false"`

    // ── Regra 2024: bônus de atributo vem do Antecedente, nunca da raça ──────
    // AbilityBonusOptions: JSON com as 3 abreviações de atributo que este
    // antecedente oferece (ex: ["INT","SAB","CAR"]), na ordem do livro. O
    // jogador distribui +2/+1 entre duas delas, ou +1 nas três (ver
    // CharacterService.Create — AbilityBonusChoice do payload de criação).
    AbilityBonusOptions string `json:"ability_bonus_options"`
    // OriginFeatName: nome do Talento (categoria "Origem") que este
    // antecedente concede automaticamente na criação — não é escolha livre
    // do jogador, é fixo por antecedente (RAW 2024, capítulo 4).
    OriginFeatName string `json:"origin_feat_name"`

    // IsLegacy: antecedente de um livro antigo (2014), sem os 3 atributos
    // nem o talento fixo que os antecedentes 2024 têm — ver a caixa lateral
    // "Antecedentes e Espécies de Livros Antigos" (PHB 2024, cap. 2, p.38):
    // o bônus de atributo vira escolha LIVRE entre os 6 atributos (mesma
    // regra +2/+1 ou +1/+1/+1), e o talento vira um talento de Origem à
    // escolha do jogador, em vez de ambos serem fixos. Não confundir com
    // simplesmente remapear pra um antecedente 2024 parecido — a regra
    // oficial é manter o antecedente antigo e liberar as duas escolhas.
    IsLegacy bool `json:"is_legacy" gorm:"default:false"`
}