package repository

import (
	"context"
	"errors"

	"github.com/guncv/Poll-Voting-Website/backend/log"
	"github.com/guncv/Poll-Voting-Website/backend/model"
	"gorm.io/gorm"
)

// UserRepository defines user database operations.
type UserRepository interface {
	CreateUser(ctx context.Context, u model.User) (model.User, error)
	FindByEmail(ctx context.Context, email string) (model.User, error)
	FindByID(ctx context.Context, id int) (model.User, error)
	UpdateUser(ctx context.Context, u model.User) (model.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db  *gorm.DB
	log log.LoggerInterface
}

// NewUserRepository creates a new userRepository with injected DB and logger.
func NewUserRepository(db *gorm.DB, logger log.LoggerInterface) UserRepository {
	return &userRepository{
		db:  db,
		log: logger,
	}
}

func (ur *userRepository) CreateUser(ctx context.Context, u model.User) (model.User, error) {
	ur.log.InfoWithID(ctx, "[Repository: CreateUser] Called for email:", u.Email)
	if err := ur.db.Create(&u).Error; err != nil {
		ur.log.ErrorWithID(ctx, "[Repository: CreateUser] Error creating user:", err)
		return model.User{}, err
	}
	ur.log.InfoWithID(ctx, "[Repository: CreateUser] Successfully created user with email:", u.Email)
	return u, nil
}

func (ur *userRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	ur.log.InfoWithID(ctx, "[Repository: FindByEmail] Called for email:", email)
	var user model.User
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.log.ErrorWithID(ctx, "[Repository: FindByEmail] User not found for email:", email)
			return model.User{}, gorm.ErrRecordNotFound
		}
		ur.log.ErrorWithID(ctx, "[Repository: FindByEmail] Error retrieving user:", err)
		return model.User{}, err
	}
	ur.log.InfoWithID(ctx, "[Repository: FindByEmail] User found for email:", email)
	return user, nil
}

func (ur *userRepository) FindByID(ctx context.Context, id int) (model.User, error) {
	ur.log.InfoWithID(ctx, "[Repository: FindByID] Called for id:", id)
	var user model.User
	if err := ur.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.log.ErrorWithID(ctx, "[Repository: FindByID] User not found with id:", id)
			return model.User{}, gorm.ErrRecordNotFound
		}
		ur.log.ErrorWithID(ctx, "[Repository: FindByID] Error finding user:", err)
		return model.User{}, err
	}
	ur.log.InfoWithID(ctx, "[Repository: FindByID] User found with id:", id)
	return user, nil
}

func (ur *userRepository) UpdateUser(ctx context.Context, u model.User) (model.User, error) {
	ur.log.InfoWithID(ctx, "[Repository: UpdateUser] Called for id:", u.UserID)
	if err := ur.db.Save(&u).Error; err != nil {
		ur.log.ErrorWithID(ctx, "[Repository: UpdateUser] Error updating user:", err)
		return model.User{}, err
	}
	ur.log.InfoWithID(ctx, "[Repository: UpdateUser] Successfully updated user with id:", u.UserID)
	return u, nil
}

func (ur *userRepository) DeleteUser(ctx context.Context, id int) error {
	ur.log.InfoWithID(ctx, "[Repository: DeleteUser] Called with id:", id)
	if err := ur.db.Delete(&model.User{}, id).Error; err != nil {
		ur.log.ErrorWithID(ctx, "[Repository: DeleteUser] Error deleting user:", err)
		return err
	}
	ur.log.InfoWithID(ctx, "[Repository: DeleteUser] Successfully deleted user with id:", id)
	return nil
}
