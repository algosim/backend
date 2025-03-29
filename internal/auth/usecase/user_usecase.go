package usecase

import (
	"github.com/algosim/backend/internal/auth/domain"
	"github.com/algosim/backend/internal/auth/repository"
	"github.com/google/uuid"
)

// UserUseCase handles user-related business logic
type UserUseCase struct {
	userRepo repository.UserRepository
}

// NewUserUseCase creates a new UserUseCase instance
func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (u *UserUseCase) CreateUser(email, oauthProvider, oauthProviderID string) (*domain.User, error) {
	// Check if user already exists
	existingUser, err := u.userRepo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Create new user
	user := domain.NewUser(email, oauthProvider, oauthProviderID)
	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (u *UserUseCase) GetUser(id uuid.UUID) (*domain.User, error) {
	return u.userRepo.FindByID(id)
}

// GetUserByEmail retrieves a user by email
func (u *UserUseCase) GetUserByEmail(email string) (*domain.User, error) {
	return u.userRepo.FindByEmail(email)
}

// GetUserByOAuthID retrieves a user by OAuth provider ID
func (u *UserUseCase) GetUserByOAuthID(provider, providerID string) (*domain.User, error) {
	return u.userRepo.FindByOAuthProviderID(provider, providerID)
}

// UpdateUser updates an existing user
func (u *UserUseCase) UpdateUser(user *domain.User) error {
	// Check if user exists
	_, err := u.userRepo.FindByID(user.ID)
	if err != nil {
		return err
	}

	// Update user
	return u.userRepo.Update(user)
}

// DeleteUser deletes a user by ID
func (u *UserUseCase) DeleteUser(id uuid.UUID) error {
	return u.userRepo.Delete(id)
}
