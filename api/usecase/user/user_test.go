package user

import (
	repo "lmm/api/domain/repository/user"
	"lmm/api/testing"
)

func TestSignup(t *testing.T) {
	tester := testing.NewTester(t)

	auth := Auth{Name: "foobar", Password: "1234"}
	requestBody := testing.StructToRequestBody(auth)

	testing.InitTable("user")
	id, err := New(repo.New()).SignUp(requestBody)
	tester.NoError(err)
	tester.Is(id, uint64(1))
}
