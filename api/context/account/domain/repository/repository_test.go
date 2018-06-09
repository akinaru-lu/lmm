package repository

import (
	"lmm/api/context/account/domain/factory"
	"lmm/api/context/account/domain/model"
	"lmm/api/db"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
)

func TestAdd(t *testing.T) {
	tester := testing.NewTester(t)
	repo := New()

	name, password := uuid.New()[:31], uuid.New()
	user, err := factory.NewUser(name, password)
	tester.NoError(err)

	err = repo.Add(user)
	tester.NoError(err)

	db := db.New()
	defer db.Close()

	stmt := db.MustPrepare("SELECT name, token FROM user WHERE id = ?")
	defer stmt.Close()

	var (
		userName  string
		userToken string
	)
	err = stmt.QueryRow(user.ID()).Scan(&userName, &userToken)

	tester.NoError(err)
	tester.Is(user.Name(), userName)
	tester.Is(user.Token(), userToken)
}

func TestFindByName_Success(t *testing.T) {
	tester := testing.NewTester(t)
	repo := New()

	name, password := uuid.New()[:31], uuid.New()
	user, _ := factory.NewUser(name, password)
	repo.Add(user)

	sameUser, err := repo.FindByName(user.Name())
	tester.NoError(err)
	tester.Isa(&model.User{}, sameUser)
	tester.Is(user.ID(), sameUser.ID())
}

func TestFindByName_NotFound(t *testing.T) {
	tester := testing.NewTester(t)
	repo := New()

	found, err := repo.FindByName("foo")

	tester.Error(err)
	tester.Nil(found)
	tester.Is(db.ErrNoRows, err)
}

func TestFindByToken_Success(t *testing.T) {
	tester := testing.NewTester(t)
	repo := New()

	name, password := uuid.New()[:31], uuid.New()
	user, _ := factory.NewUser(name, password)
	repo.Add(user)

	found, err := repo.FindByToken(user.Token())
	tester.NoError(err)
	tester.Isa(&model.User{}, found)
	tester.Is(user.ID(), found.ID())
}

func TestFindByToken_NotFound(t *testing.T) {
	tester := testing.NewTester(t)
	repo := New()

	found, err := repo.FindByToken("1234")
	tester.Error(err)
	tester.Is(db.ErrNoRows, err)
	tester.Nil(found)
}
