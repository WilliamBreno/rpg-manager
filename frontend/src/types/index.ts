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
  is_default: boolean
}

export interface Race {
  ID: number
  name: string
  edition: string
  description: string
  speed: number
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

  // Palavras-chave ex: "Arcano, Arma"
  keywords?: string

  // Mecânicas
  action_type?: string
  range?: string
  target?: string
  attack?: string

  // Efeitos
  description?: string
  hit?: string
  miss?: string
  effect?: string
  special?: string
  level_scaling?: string

  // Características de classe
  is_class_feature?: boolean
  requires_choice?: boolean
  choice_group?: string
  is_race_feature: boolean
}

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