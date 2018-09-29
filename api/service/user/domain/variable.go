package domain

import "github.com/pkg/errors"

var (
	// ErrInvalidUserID error
	ErrInvalidUserID = errors.New("invalid user id")

	// ErrInvalidUserName error
	ErrInvalidUserName = errors.New("invalid user name")

	// ErrInvalidUserToken error
	ErrInvalidUserToken = errors.New("invalid user token, expect a uuid v4")

	// ErrNoSuchUser error
	ErrNoSuchUser = errors.New("no such user")

	// ErrUserNameAlreadyUsed error
	ErrUserNameAlreadyUsed = errors.New("user name has already been used")

	// ErrUserPasswordEmpty error
	ErrUserPasswordEmpty = errors.New("user password is empty")

	// ErrUserPasswordTooShort error
	ErrUserPasswordTooShort = errors.New("user password should be equal to or longer than 8")

	// ErrUserPasswordTooWeak error
	ErrUserPasswordTooWeak = errors.New("user password is too weak")
)
