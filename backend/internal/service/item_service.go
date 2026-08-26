package service

import (
	"rpg-manager/internal/domain"
	"rpg-manager/internal/repository"
)

type ItemService struct{ Repo *repository.ItemRepository }

func NewItemService(repo *repository.ItemRepository) *ItemService { return &ItemService{Repo: repo} }

func (s *ItemService) GetAll(edition, category string) ([]domain.Item, error) {
	return s.Repo.FindAll(edition, category)
}

func (s *ItemService) GetByID(id uint) (domain.Item, error) {
	return s.Repo.FindByID(id)
}
