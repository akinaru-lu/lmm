package usecase

import (
	"lmm/api/context/account/domain/repository"
	testingService "lmm/api/context/account/domain/service/testing"
	"lmm/api/testing"
)

func TestSignUp(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)

	auth := Auth{Name: "foobar", Password: "1234"}
	id, err := New(repository.New()).SignUp(auth.Name, auth.Password)
	tester.NoError(err)
	tester.Is(uint64(1), id)
}

func TestSignUp_Duplicate(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)

	user := testingService.NewUser()

	repo := repository.New()
	id, err := New(repo).SignUp(user.Name, user.Password)
	tester.Error(err)
	tester.Is(ErrDuplicateUserName.Error(), err.Error())
	tester.Is(uint64(0), id)
}
