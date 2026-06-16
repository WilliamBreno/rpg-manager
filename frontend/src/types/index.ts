export interface Class {
  ID: number
  name: string
  edition: string
  description: string

  // 5e
  hit_die: number

  // 4e
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
  description?: string    // opcional — nem sempre preenchido
  power_type?: PowerType  // opcional — habilidades 5e podem não ter tipo
  level?: number
  edition?: string
  class_id?: number | null
  race_id?: number | null
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

  // HP
  hit_points: number
  max_hp: number
  temp_hp: number

  // Pulsos de Cura (4e)
  surge_value: number
  surges_per_day: number

  // Defesas calculadas (4e)
  defense_ac: number
  defense_fort: number
  defense_refl: number
  defense_will: number

  // Relacionamentos
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

  // Objetos relacionados
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