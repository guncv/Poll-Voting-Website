package service

import (
	"context"
	"errors"

	"github.com/guncv/Poll-Voting-Website/backend/log"
	"github.com/guncv/Poll-Voting-Website/backend/model"
	"github.com/guncv/Poll-Voting-Website/backend/repository"
	"github.com/guncv/Poll-Voting-Website/backend/util"
	"gorm.io/gorm"
)

// UserService defines the methods for user operations.
type UserService interface {
	Register(ctx context.Context, email, password string) (model.User, error)
	Login(ctx context.Context, email, password string) (model.User, error)
	GetUserByID(ctx context.Context, id string) (model.User, error)
	UpdateUser(ctx context.Context, id string, newEmail, newPassword string) (model.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	repo                repository.UserRepository
	log                 log.LoggerInterface
	notificationService INotificationService
}

// NewUserService creates a new userService with injected repository and logger.
func NewUserService(r repository.UserRepository, logger log.LoggerInterface, notificationService INotificationService) UserService {
	return &userService{
		repo:                r,
		log:                 logger,
		notificationService: notificationService,
	}
}

// Register creates a new user if the email is not already taken.
func (us *userService) Register(ctx context.Context, email, password string) (model.User, error) {
	us.log.InfoWithID(ctx, "[Service: Register] Called with email:", email)

	// Check if the user already exists.
	_, err := us.repo.FindByEmail(ctx, email)
	if err == nil {
		us.log.ErrorWithID(ctx, "[Service: Register] User already exists for email:", email)
		return model.User{}, errors.New("user already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		us.log.ErrorWithID(ctx, "[Service: Register] Error checking existing user:", err)
		return model.User{}, err
	}

	// Hash the password before storing it.
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		us.log.ErrorWithID(ctx, "[Service: Register] Error hashing password:", err)
		return model.User{}, err
	}

	newUser := model.User{
		Email:    email,
		Password: hashedPassword,
	}
	created, err := us.repo.CreateUser(ctx, newUser)
	if err != nil {
		us.log.ErrorWithID(ctx, "[Service: Register] Error creating user:", err)
		return model.User{}, err
	}

	err = us.notificationService.AddSubscriberToUserTopic(ctx, email)
	if err != nil {
		us.log.ErrorWithID(ctx, "[Service: Register] Error adding subscriber to admin topic:", err)
		return model.User{}, err
	}

	us.log.InfoWithID(ctx, "[Service: Register] User created successfully with email:", email)
	return created, nil
}

// Login checks the user's credentials.
func (us *userService) Login(ctx context.Context, email, password string) (model.User, error) {
	us.log.InfoWithID(ctx, "[Service: Login] Called with email:", email)

	u, err := us.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			us.log.ErrorWithID(ctx, "[Service: Login] User not found for email:", email)
			return model.User{}, errors.New("invalid credentials")
		}
		us.log.ErrorWithID(ctx, "[Service: Login] Error retrieving user:", err)
		return model.User{}, err
	}

	// Verify the provided password against the stored hash.
	if err := util.CheckPassword(password, u.Password); err != nil {
		us.log.ErrorWithID(ctx, "[Service: Login] Invalid credentials for email:", email)
		return model.User{}, errors.New("invalid credentials")
	}

	err = us.notificationService.NotifyUserOfUrgentQuestion(ctx, "Urgent Question", "A new user has logged in with email: "+email)
	if err != nil {
		us.log.ErrorWithID(ctx, "[Service: Login] Error notifying user of urgent question:", err)
		return model.User{}, err
	}

	us.log.InfoWithID(ctx, "[Service: Login] User logged in successfully with email:", email)
	return u, nil
}

// GetUserByID retrieves a user by their ID.
func (us *userService) GetUserByID(ctx context.Context, id string) (model.User, error) {
	us.log.InfoWithID(ctx, "[Service: GetUserByID] Called with id:", id)

	u, err := us.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			us.log.ErrorWithID(ctx, "[Service: GetUserByID] User not found with id:", id)
			return model.User{}, errors.New("user not found")
		}
		us.log.ErrorWithID(ctx, "[Service: GetUserByID] Error finding user:", err)
		return model.User{}, err
	}

	us.log.InfoWithID(ctx, "[Service: GetUserByID] User found with id:", id)
	return u, nil
}

// UpdateUser modifies an existing user's email and/or password.
func (us *userService) UpdateUser(ctx context.Context, id string, newEmail, newPassword string) (model.User, error) {
	us.log.InfoWithID(ctx, "[Service: UpdateUser] Called with id:", id)

	u, err := us.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			us.log.ErrorWithID(ctx, "[Service: UpdateUser] User not found with id:", id)
			return model.User{}, errors.New("user not found")
		}
		us.log.ErrorWithID(ctx, "[Service: UpdateUser] Error finding user:", err)
		return model.User{}, err
	}

	if newEmail != "" {
		u.Email = newEmail
	}
	if newPassword != "" {
		// Hash the new password.
		hashedPassword, err := util.HashPassword(newPassword)
		if err != nil {
			us.log.ErrorWithID(ctx, "[Service: UpdateUser] Error hashing new password:", err)
			return model.User{}, err
		}
		u.Password = hashedPassword
	}

	updated, err := us.repo.UpdateUser(ctx, u)
	if err != nil {
		us.log.ErrorWithID(ctx, "[Service: UpdateUser] Error updating user:", err)
		return model.User{}, err
	}

	us.log.InfoWithID(ctx, "[Service: UpdateUser] User updated successfully with id:", id)
	return updated, nil
}

// DeleteUser ensures a user exists and deletes the user.
func (us *userService) DeleteUser(ctx context.Context, id string) error {
	us.log.InfoWithID(ctx, "[Service: DeleteUser] Called with id:", id)

	// Check if the user exists.
	_, err := us.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			us.log.ErrorWithID(ctx, "[Service: DeleteUser] User not found")
			return errors.New("user not found")
		}
		us.log.ErrorWithID(ctx, "[Service: DeleteUser] Error finding user:", err)
		return err
	}

	// Proceed to delete the user.
	err = us.repo.DeleteUser(ctx, id)
	if err != nil {
		us.log.ErrorWithID(ctx, "[Service: DeleteUser] Error deleting user:", err)
		return err
	}

	us.log.InfoWithID(ctx, "[Service: DeleteUser] Successfully deleted user with id:", id)
	return nil
}
