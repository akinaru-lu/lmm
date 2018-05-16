package usecase

import (
	"lmm/api/context/account/domain/model"
	"lmm/api/domain/repository"
)

func (uc *Usecase) SignUp(name, password string) (uint64, error) {
	if name == "" || password == "" {
		return 0, ErrEmptyUserNameOrPassword
	}

	m := model.New(name, password)
	user, err := uc.repo.Save(m)
	if err != nil {
		key, _, ok := repository.CheckErrorDuplicate(err.Error())
		if !ok {
			return 0, err
		}
		if key == "name" {
			return 0, ErrDuplicateUserName
		}
		return 0, err
	}
	return user.ID, nil
}
