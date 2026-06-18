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
  // Perícias (4e)
  trained_skills_count: number  // quantas pode treinar
  available_skills: string      // JSON: '["Atletismo","Percepção"]'
  // Talentos (4e)
  talentos_count: number        // 2 na campanha
}

export interface Race {
  ID: number
  name: string
  edition: string
  description: string
  speed: number
  // Perícias (4e)
  bonus_trained_skills: number  // 0 para maioria, 1 para Humano/Meio-Elfo
  bonus_skill_values: string    // JSON: '{"Percepção": 2}'
  // Talentos (4e)
  bonus_talentos: number        // 0 para maioria, 1 para Humano
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

// ── NOVO — Perícia ──────────────────────────────────────────────────────────
export interface Pericia {
  ID: number
  name: string
  attribute: string   // "Força", "Destreza", etc.
  description: string
  tooltip: string
  edition: string
}

export interface CharacterPericia {
  character_id: number
  pericia_name: string
}

// ── NOVO — Talento ──────────────────────────────────────────────────────────
export interface Talento {
  ID: number
  name: string
  edition: string
  category: string      // "Combate", "Defesa", "Perícia", "Magia", "Armadura"
  description: string
  prerequisite: string  // "" se não tiver
  tooltip: string
}

// ── Background ──────────────────────────────────────────────────────────────
export interface Background {
  ID?: number
  character_id?: number
  history: string
  personality_traits: string
  ideals: string
  bonds: string
  flaws: string
}

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
  defense_ac: number
  defense_fort: number
  defense_refl: number
  defense_will: number
  armor_id?: number
  class_id: number
  race_id: number
  avatar_url?: string
  strength: number
  dexterity: number
  constitution: number
  intelligence: number
  wisdom: number
  charisma: number
  class: Class
  race: Race
  armor?: Armor
  skills: Skill[]
  background?: Background
  // NOVO
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
}