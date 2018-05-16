package ui

import (
	"io"
	"lmm/api/context/account/domain/model"
	"lmm/api/context/account/domain/repository"
	"lmm/api/context/account/usecase"
	"lmm/api/testing"
	"net/http"
)

func TestPostV1Signup(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)

	router := testing.NewRouter()
	router.POST("/v1/signup", SignUp)

	reqeustBody := testing.StructToRequestBody(usecase.Auth{Name: "foobar", Password: "1234"})

	res := testing.NewResponse()
	router.ServeHTTP(res, testing.POST("/v1/signup", reqeustBody))

	tester.Is(http.StatusCreated, res.StatusCode())
	tester.Is("/users/1", res.Header().Get("Location"))
}

func TestPostV1Signup_Duplicate(t *testing.T) {
	testing.InitTable("user")
	tester := testing.NewTester(t)

	router := testing.NewRouter()
	router.POST("/v1/signup", SignUp)

	model := model.New("foobar", "1234")
	repository.New().Save(model)

	res := testing.NewResponse()
	auth := usecase.Auth{Name: "foobar", Password: "1234"}
	router.ServeHTTP(res, testing.POST("/v1/signup", testing.StructToRequestBody(auth)))

	tester.Is(http.StatusBadRequest, res.StatusCode())
	tester.Is(usecase.ErrDuplicateUserName.Error()+"\n", res.Body())
}

func TestPostV1SignUp_400_EmptyUserName(t *testing.T) {
	testing.InitTable("user")

	requestBody := testing.StructToRequestBody(Auth{Name: "", Password: "1234"})
	res := postSignUp(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusBadRequest, res.StatusCode())
	tester.Is(usecase.ErrEmptyUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignUp_400_EmptyPassword(t *testing.T) {
	testing.InitTable("user")

	requestBody := testing.StructToRequestBody(Auth{Name: "foobar", Password: ""})
	res := postSignUp(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusBadRequest, res.StatusCode())
	tester.Is(usecase.ErrEmptyUserNameOrPassword.Error()+"\n", res.Body())
}

func TestPostV1SignUp_400_EmptyUserNameAndPassword(t *testing.T) {
	testing.InitTable("user")

	requestBody := testing.StructToRequestBody(Auth{Name: "", Password: ""})
	res := postSignUp(requestBody)

	tester := testing.NewTester(t)
	tester.Is(http.StatusBadRequest, res.StatusCode())
	tester.Is(usecase.ErrEmptyUserNameOrPassword.Error()+"\n", res.Body())
}

func postSignUp(requestBody io.Reader) *testing.Response {
	res := testing.NewResponse()

	router := testing.NewRouter()
	router.POST("/v1/signup", SignIn)
	router.ServeHTTP(res, testing.POST("/v1/signup", requestBody))

	return res
}
