package model

import (
	"regexp"
	"time"

	"github.com/pkg/errors"

	"lmm/api/clock"
	"lmm/api/model"
	"lmm/api/service/user/domain"
	"lmm/api/util/uuidutil"
)

var (
	patternUserName = regexp.MustCompile(`^[a-zA-Z]{1}[0-9a-zA-Z_-]{2,17}$`)
)

// UserDescriptor describes user's basic infomation
type UserDescriptor struct {
	model.Entity
	name         string
	role         Role
	registeredAt time.Time
}

// NewUserDescriptor creates a new *UserDescriptor
func NewUserDescriptor(name string, role Role, registeredAt time.Time) (*UserDescriptor, error) {
	user := UserDescriptor{
		role:         role,
		registeredAt: registeredAt,
	}

	if err := user.setName(name); err != nil {
		return nil, err
	}

	return &user, nil
}

// Name gets user's name
func (user *UserDescriptor) Name() string {
	return user.name
}

// Is compares if two use are the same
func (user *UserDescriptor) Is(target *UserDescriptor) bool {
	return user.Name() == target.Name()
}

// Role gets user's Role model
func (user *UserDescriptor) Role() Role {
	return user.role
}

// RegisteredAt gets user's register date
func (user *UserDescriptor) RegisteredAt() time.Time {
	return user.registeredAt
}

func (user *UserDescriptor) setName(name string) error {
	if !patternUserName.MatchString(name) {
		return domain.ErrInvalidUserName
	}
	user.name = name
	return nil
}

// User domain model
type User struct {
	UserDescriptor
	password string
	token    string
}

// NewUser creates a new user domain model
func NewUser(name, password, token string, role Role) (*User, error) {
	user := User{}

	if err := user.setName(name); err != nil {
		return nil, err
	}

	if err := user.setPassword(password); err != nil {
		return nil, err
	}

	if err := user.setToken(token); err != nil {
		return nil, err
	}

	if err := user.ChangeRole(role); err != nil {
		return nil, err
	}

	user.registeredAt = clock.Now()

	return &user, nil
}

// ChangeRole changes user's role
func (user *User) ChangeRole(role Role) error {
	switch role {
	case Admin, Guest, Ordinary:
		user.role = role
		return nil
	default:
		return domain.ErrNoSuchRole
	}
}

// Password gets user's encrypted password
func (user *User) Password() string {
	return user.password
}

// ChangePassword changes password to given newPassword
func (user *User) ChangePassword(newPassword string) error {
	return user.setPassword(newPassword)
}

func (user *User) setPassword(password string) error {
	user.password = password
	return nil
}

// Token gets user's token
func (user *User) Token() string {
	return user.token
}

func (user *User) setToken(token string) error {
	uuid, err := uuidutil.ParseString(token)
	if err != nil {
		return errors.Wrap(domain.ErrInvalidUserToken, err.Error())
	}
	if v := uuid.Version().String(); v != "VERSION_4" {
		return errors.Wrap(domain.ErrInvalidUserToken, "unexpected uuid version: "+v)
	}
	user.token = token
	return nil
}

// Is compares if two users are the same
func (user *User) Is(other *User) bool {
	return user.UserDescriptor.Is(&other.UserDescriptor)
}
