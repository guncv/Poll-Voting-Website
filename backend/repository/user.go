package repository

import (
	"errors"

	"github.com/guncv/Poll-Voting-Website/backend/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(u model.User) (model.User, error)
	FindByEmail(email string) (model.User, error)
	FindByID(id int) (model.User, error)
	UpdateUser(u model.User) (model.User, error)
	DeleteUser(id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) CreateUser(u model.User) (model.User, error) {
	if err := ur.db.Create(&u).Error; err != nil {
		return model.User{}, err
	}
	return u, nil
}

func (ur *userRepository) FindByEmail(email string) (model.User, error) {
	var user model.User
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, gorm.ErrRecordNotFound
		}
		return model.User{}, err
	}
	return user, nil
}

func (ur *userRepository) FindByID(id int) (model.User, error) {
	var user model.User
	if err := ur.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, gorm.ErrRecordNotFound
		}
		return model.User{}, err
	}
	return user, nil
}

func (ur *userRepository) UpdateUser(u model.User) (model.User, error) {
	if err := ur.db.Save(&u).Error; err != nil {
		return model.User{}, err
	}
	return u, nil
}

func (ur *userRepository) DeleteUser(id int) error {
	if err := ur.db.Delete(&model.User{}, id).Error; err != nil {
		return err
	}
	return nil
}
