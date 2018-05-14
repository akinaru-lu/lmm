package usecase

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/service"
	"lmm/api/db"
	"lmm/api/utils/sha256"
)

func (uc *Usecase) SignIn(name, password string) (*model.User, error) {
	if name == "" || password == "" {
		return nil, ErrEmptyUserNameOrPassword
	}

	user, err := uc.repo.FindByName(name)
	if err != nil {
		if err == db.ErrNoRows {
			return nil, ErrInvalidUserNameOrPassword
		}
		return nil, err
	}

	// validate password
	encoded := sha256.Hex([]byte(user.GUID + password))
	if encoded != user.Password {
		return nil, ErrInvalidUserNameOrPassword
	}

	encodedToken := service.EncodeToken(user.Token)
	user.Token = encodedToken

	return user, nil
}
