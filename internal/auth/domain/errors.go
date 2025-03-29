package domain

import "errors"

var (
	// ErrUserNotFound is returned when a user cannot be found
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists is returned when trying to create a user that already exists
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrInvalidCredentials is returned when login credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrTokenNotFound is returned when a token cannot be found
	ErrTokenNotFound = errors.New("token not found")

	// ErrTokenExpired is returned when a token has expired
	ErrTokenExpired = errors.New("token expired")

	// ErrInvalidToken is returned when a token is invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrOAuthProviderNotSupported is returned when the OAuth provider is not supported
	ErrOAuthProviderNotSupported = errors.New("oauth provider not supported")

	// ErrOAuthCallbackFailed is returned when OAuth callback fails
	ErrOAuthCallbackFailed = errors.New("oauth callback failed")
)
