export interface Class {
  ID: number
  name: string
  edition: string
  description: string
  hit_die: number
  base_hp: number
  hp_per_level: number
  surges_per_day: number
  fort_bonus: number
  refl_bonus: number
  will_bonus: number
  saving_throws?: string
  is_default: boolean
  // Perícias
  trained_skills_count: number
  available_skills: string       // JSON: '["Atletismo","Percepção"]'
  // Talentos (4e)
  talentos_count: number
  automatic_pericias?: string
}

export interface Race {
  ID: number
  name: string
  edition: string
  description: string
  speed: number
  // Perícias (4e)
  bonus_trained_skills: number
  bonus_skill_values: string     // JSON: '{"Percepção": 2}'
  // Talentos (4e)
  bonus_talentos: number
}

export type PowerType = 'utility' | 'unlimited' | 'encounter' | 'daily'

export interface Skill {
  ID: number
  name: string
  edition?: string
  class_id?: number | null
  race_id?: number | null
  power_type?: PowerType
  level?: number
  keywords?: string
  action_type?: string
  range?: string
  target?: string
  attack?: string
  description?: string
  hit?: string
  miss?: string
  effect?: string
  special?: string
  level_scaling?: string
  is_class_feature?: boolean
  requires_choice?: boolean
  choice_group?: string
  is_race_feature: boolean
}

// ── Perícia ─────────────────────────────────────────────────────────────────
export interface Pericia {
  ID: number
  name: string
  attribute: string
  description: string
  tooltip: string
  edition: string
}

export interface CharacterPericia {
  character_id: number
  pericia_name: string
}

// ── Talento ─────────────────────────────────────────────────────────────────
export interface Talento {
  ID: number
  name: string
  edition: string
  category: string
  description: string
  prerequisite: string
  tooltip: string
}

// ── Background — biografia/personalidade do personagem (sistema existente) ──
export interface Background {
  ID?: number
  character_id?: number
  history: string
  personality_traits: string
  ideals: string
  bonds: string
  flaws: string
}

// ── Antecedent — Antecedente D&D 5e (Acólito, Criminoso, Soldado, etc.) ────
export interface Antecedent {
  ID: number
  name: string
  edition: string
  description: string
  skill_proficiencies: string    // JSON: '["Intuição","Religião"]'
  tool_proficiencies: string
  languages: string
  equipment: string
  feature: string
  feature_description: string
  is_default: boolean
}

// ── Armor ────────────────────────────────────────────────────────────────────
export type ArmorType = 'none' | 'light' | 'medium' | 'heavy' | 'shield'

export interface Armor {
  ID: number
  name: string
  edition: string
  armor_type: ArmorType
  base_ac: number
  max_dex_bonus: number
  description: string
}

// ── Character ────────────────────────────────────────────────────────────────
export interface Character {
  ID: number
  name: string
  edition: string
  level: number
  hit_points: number
  max_hp: number
  temp_hp: number
  surge_value: number
  surges_per_day: number
  // Defesas (4e)
  defense_ac: number
  defense_fort: number
  defense_refl: number
  defense_will: number
  armor_id?: number
  class_id: number
  race_id: number
  avatar_url?: string
  // Atributos
  strength: number
  dexterity: number
  constitution: number
  intelligence: number
  wisdom: number
  charisma: number
  // ── Campos D&D 5e ──────────────────────────────────────────────────────────
  antecedent_id?: number          // FK para Antecedent
  alignment?: string              // "Leal e Bom", "Neutro", etc.
  personality_traits?: string
  ideals?: string
  bonds?: string
  flaws?: string
  speed?: number                  // deslocamento em pés (30 para maioria)
  proficiency_bonus?: number      // calculado: +2 nos níveis 1-4
  // ── Relacionamentos ────────────────────────────────────────────────────────
  class: Class
  race: Race
  armor?: Armor
  background?: Background         // biografia/notas (sistema existente)
  antecedent?: Antecedent         // Antecedente D&D 5e
  skills: Skill[]
  pericias: CharacterPericia[]
  talentos: Talento[]
}

export interface CreateCharacterDTO {
  name: string
  edition: string
  class_id: number
  race_id: number
  hit_points: number
  strength: number
  dexterity: number
  constitution: number
  intelligence: number
  wisdom: number
  charisma: number
  armor_id?: number
  // 5e
  antecedent_id?: number
  alignment?: string
  personality_traits?: string
  ideals?: string
  bonds?: string
  flaws?: string
}
export interface Antecedent {
  ID: number
  name: string
  edition: string
  description: string
  skill_proficiencies: string
  tool_proficiencies: string
  languages: string
  equipment: string
  feature: string
  feature_description: string
  is_default: boolean
  CreatedAt?: string
  UpdatedAt?: string
}