package repository

import (
    "rpg-manager/internal/domain"
    "gorm.io/gorm"
)

type SkillRepository struct {
    DB *gorm.DB
}

func NewSkillRepository(db *gorm.DB) *SkillRepository {
    return &SkillRepository{DB: db}
}

func (r *SkillRepository) FindAll() ([]domain.Skill, error) {
    var skills []domain.Skill
    result := r.DB.Find(&skills)
    return skills, result.Error
}

func (r *SkillRepository) FindByID(id uint) (domain.Skill, error) {
    var skill domain.Skill
    result := r.DB.First(&skill, id)
    return skill, result.Error
}

func (r *SkillRepository) FindByClassID(classID uint) ([]domain.Skill, error) {
    var skills []domain.Skill
    result := r.DB.Where("class_id = ?", classID).Find(&skills)
    return skills, result.Error
}

func (r *SkillRepository) FindByRaceID(raceID uint) ([]domain.Skill, error) {
    var skills []domain.Skill
    result := r.DB.Where("race_id = ?", raceID).Find(&skills)
    return skills, result.Error
}

func (r *SkillRepository) FindByClassAndRace(classID, raceID uint) ([]domain.Skill, error) {
    var skills []domain.Skill
    result := r.DB.Where("class_id = ? OR race_id = ?", classID, raceID).Find(&skills)
    return skills, result.Error
}

func (r *SkillRepository) FindByLevel(classID uint, level int) ([]domain.Skill, error) {
    var skills []domain.Skill
    result := r.DB.Where("class_id = ? AND level <= ?", classID, level).Find(&skills)
    return skills, result.Error
}

func (r *SkillRepository) Create(skill *domain.Skill) error {
    return r.DB.Create(skill).Error
}

func (r *SkillRepository) Update(skill *domain.Skill) error {
    return r.DB.Save(skill).Error
}

func (r *SkillRepository) Delete(id uint) error {
    return r.DB.Delete(&domain.Skill{}, id).Error
}