package service

import (
    "errors"

    "github.com/guncv/Poll-Voting-Website/backend/model"
    "github.com/guncv/Poll-Voting-Website/backend/repository"
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

// NewUserService injects the UserRepository into the service.
func NewUserService(r repository.UserRepository) UserService {
    return &userService{repo: r}
}

// Register creates a new user if the email doesn't already exist.
func (us *userService) Register(email, password string) (model.User, error) {
    // Check if user already exists
	_, err := us.repo.FindByEmail(email)
	if err == nil {
		return model.User{}, errors.New("user already exists")
	}
	
    if !errors.Is(err, gorm.ErrRecordNotFound) {
        // DB error (connection issue, etc.)
        return model.User{}, err
    }

    // Create user
    newUser := model.User{
        Email:    email,
        Password: password, // In production, always hash the password
    }
    created, err := us.repo.CreateUser(newUser)
    if err != nil {
        return model.User{}, err
    }
    return created, nil
}

// Login checks if the email/password match an existing user.
func (us *userService) Login(email, password string) (model.User, error) {
    user, err := us.repo.FindByEmail(email)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return model.User{}, errors.New("invalid credentials")
        }
        return model.User{}, err
    }

    // Check password match (in real app, compare hash)
    if user.Password != password {
        return model.User{}, errors.New("invalid credentials")
    }

    return user, nil
}

// GetUserByID retrieves a user by ID.
func (us *userService) GetUserByID(id int) (model.User, error) {
    user, err := us.repo.FindByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return model.User{}, errors.New("user not found")
        }
        return model.User{}, err
    }
    return user, nil
}

// UpdateUser modifies the user’s email or password if provided.
func (us *userService) UpdateUser(id int, newEmail, newPassword string) (model.User, error) {
    // First fetch existing user
    user, err := us.repo.FindByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return model.User{}, errors.New("user not found")
        }
        return model.User{}, err
    }

    if newEmail != "" {
        user.Email = newEmail
    }
    if newPassword != "" {
        user.Password = newPassword // again, hash it in real usage
    }

    updated, err := us.repo.UpdateUser(user)
    if err != nil {
        return model.User{}, err
    }
    return updated, nil
}

// DeleteUser removes a user by ID.
func (us *userService) DeleteUser(id int) error {
    // Check if user exists first
    _, err := us.repo.FindByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("user not found")
        }
        return err
    }

    // Then delete
    if err := us.repo.DeleteUser(id); err != nil {
        return err
    }
    return nil
}
