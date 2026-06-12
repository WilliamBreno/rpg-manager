package repository

import (
    "rpg-manager/internal/domain"
    "gorm.io/gorm"
)

type UserRepository struct {
    DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{DB: db}
}

func (r *UserRepository) FindByEmail(email string) (domain.User, error) {
    var user domain.User
    result := r.DB.Where("email = ?", email).First(&user)
    return user, result.Error
}

func (r *UserRepository) FindByID(id uint) (domain.User, error) {
    var user domain.User
    result := r.DB.First(&user, id)
    return user, result.Error
}

func (r *UserRepository) Create(user *domain.User) error {
    return r.DB.Create(user).Error
}