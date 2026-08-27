package service

import (
	"rpg-manager/internal/domain"
	"rpg-manager/internal/repository"
)

type ClassEquipmentService struct{ Repo *repository.ClassEquipmentRepository }

func NewClassEquipmentService(repo *repository.ClassEquipmentRepository) *ClassEquipmentService {
	return &ClassEquipmentService{Repo: repo}
}

func (s *ClassEquipmentService) GetByClass(classID uint) ([]domain.ClassEquipmentOption, error) {
	return s.Repo.GetByClass(classID)
}
