package repository

import (
    "errors"

    "github.com/guncv/Poll-Voting-Website/backend/model"
    "gorm.io/gorm"
)

// UserRepository defines the data-access methods for the 'users' table.
type UserRepository interface {
    CreateUser(user model.User) (model.User, error)
    FindByEmail(email string) (model.User, error)
    FindByID(id int) (model.User, error)
    UpdateUser(user model.User) (model.User, error)
    DeleteUser(id int) error
}

type userRepository struct {
    db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

// CreateUser inserts a new user record in the DB.
func (ur *userRepository) CreateUser(u model.User) (model.User, error) {
    if err := ur.db.Create(&u).Error; err != nil {
        return model.User{}, err
    }
    return u, nil
}

// FindByEmail searches for a user by email.
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

// FindByID searches for a user by ID.
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

// UpdateUser saves changes to an existing user.
func (ur *userRepository) UpdateUser(u model.User) (model.User, error) {
    if err := ur.db.Save(&u).Error; err != nil {
        return model.User{}, err
    }
    return u, nil
}

// DeleteUser removes a user by ID.
func (ur *userRepository) DeleteUser(id int) error {
    // GORM doesn't return an error for zero rows affected. 
    // If you care, you can check RowsAffected after the operation.
    if err := ur.db.Delete(&model.User{}, id).Error; err != nil {
        return err
    }
    return nil
}
