package service

import (
	"errors"

	"github.com/guncv/Poll-Voting-Website/backend/model"
	"github.com/guncv/Poll-Voting-Website/backend/repository"
	"github.com/guncv/Poll-Voting-Website/backend/util"
	"gorm.io/gorm"
)

type UserService interface {
	Register(email, password string) (model.User, error)
	Login(email, password string) (model.User, error)
	GetUserByID(id int) (model.User, error)
	UpdateUser(id int, newEmail, newPassword string) (model.User, error)
	DeleteUser(id int) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

// Register a new user if email not taken
func (us *userService) Register(email, password string) (model.User, error) {
	// Check if user already exists
	_, err := us.repo.FindByEmail(email)
	if err == nil {
		return model.User{}, errors.New("user already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, err
	}

	// Hash the password before storing it
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return model.User{}, err
	}

	newUser := model.User{
		Email:    email,
		Password: hashedPassword, // store hashed password
	}
	created, err := us.repo.CreateUser(newUser)
	if err != nil {
		return model.User{}, err
	}
	return created, nil
}

// Login checks email/password
func (us *userService) Login(email, password string) (model.User, error) {
	u, err := us.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, errors.New("invalid credentials")
		}
		return model.User{}, err
	}

	// Check if the provided password matches the stored hashed password
	if err := util.CheckPassword(password, u.Password); err != nil {
		return model.User{}, errors.New("invalid credentials")
	}
	return u, nil
}


// GetUserByID returns user or "user not found"
func (us *userService) GetUserByID(id int) (model.User, error) {
	u, err := us.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}
	return u, nil
}

// UpdateUser modifies user if found
func (us *userService) UpdateUser(id int, newEmail, newPassword string) (model.User, error) {
	u, err := us.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}

	if newEmail != "" {
		u.Email = newEmail
	}
	if newPassword != "" {
		u.Password = newPassword // normally hashed
	}

	updated, err := us.repo.UpdateUser(u)
	if err != nil {
		return model.User{}, err
	}
	return updated, nil
}

// DeleteUser ensures user exists, then deletes
func (us *userService) DeleteUser(id int) error {
	_, err := us.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}
	return us.repo.DeleteUser(id)
}
