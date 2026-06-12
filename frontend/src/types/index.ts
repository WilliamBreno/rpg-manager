export interface Class {
  ID: number
  name: string
  edition: string
  description: string
  hit_die: number
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
  description: string
  power_type: PowerType
  level: number
  edition: string
  class_id: number | null
  race_id: number | null
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

export interface Character {
  ID: number
  name: string
  edition: string
  level: number
  hit_points: number
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