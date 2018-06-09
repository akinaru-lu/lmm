package ui

import (
	"encoding/json"
	"lmm/api/context/account/domain/factory"
	"lmm/api/context/account/domain/repository"
	"lmm/api/context/account/domain/service"
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/usecase/auth"
	"lmm/api/utils/uuid"
	"strings"
)

func TestGETV1Verify_Success(t *testing.T) {
	tester := testing.NewTester(t)

	name, password := uuid.New()[:31], uuid.New()
	user, _ := factory.NewUser(name, password)
	repository.New().Add(user)

	headers := make(map[string]string, 0)
	headers["Authorization"] = "Bearer " + service.EncodeToken(user.Token())
	res := getVerify(headers)

	schema := VerifyResponse{}
	json.NewDecoder(strings.NewReader(res.Body())).Decode(&schema)

	tester.Is(http.StatusOK, res.StatusCode())
	tester.Is(user.ID(), schema.ID)
	tester.Is(user.Name(), schema.Name)
}

func TestGETV1Verify_NoAuthorization(t *testing.T) {
	tester := testing.NewTester(t)

	res := getVerify(nil)

	tester.Is(http.StatusUnauthorized, res.StatusCode())
}

func TestGETV1Verify_NotBearerAuthorization(t *testing.T) {
	tester := testing.NewTester(t)

	name, password := uuid.New()[:31], uuid.New()
	user, _ := factory.NewUser(name, password)
	repository.New().Add(user)

	headers := make(map[string]string, 0)
	headers["Authorization"] = service.EncodeToken(user.Token())
	res := getVerify(headers)

	tester.Is(http.StatusUnauthorized, res.StatusCode())
}

func TestGETV1Verify_InvalidToken(t *testing.T) {
	tester := testing.NewTester(t)

	headers := make(map[string]string, 0)
	headers["Authorization"] = "Bearer xxx"
	res := getVerify(headers)

	tester.Is(http.StatusUnauthorized, res.StatusCode())
}

func getVerify(headers map[string]string) *testing.Response {
	request := testing.GET("/v1/verify")
	if headers != nil {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}

	res := testing.NewResponse()

	router := testing.NewRouter()
	router.GET("/v1/verify", auth.BearerAuth(Verify))
	router.ServeHTTP(res, request)

	return res
}
