package testingutil

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"lmm/api/service/auth/domain/model"
	"lmm/api/storage/db"
	"lmm/api/util/stringutil"
)

// NewAuthUser creates new user from auth service
func NewAuthUser(db db.DB) (*model.User, error) {
	rawPassword := uuid.New().String()
	b, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	encryptedPassword := string(b)

	user := model.NewUser(
		uuid.New().String()[:8],
		encryptedPassword,
		stringutil.ReplaceAll(uuid.New().String(), "-", ""),
	)

	now := time.Now()

	if _, err := db.Exec(context.Background(), `
		insert into user (name, password, token, created_at) values (?, ?, ?, ?)
	`, user.Name(), encryptedPassword, user.RawToken(), now); err != nil {
		panic(err)
	}

	return user, nil
}

// ExtractAccessToken tries to extract access token from given string
func ExtractAccessToken(s string) (string, error) {
	// avoid cycle import, see lmm/api/service/auth/ui/adapter.go
	type loginResponse struct {
		AccessToken string `json:"accessToken"`
	}

	schema := loginResponse{}

	if err := json.NewDecoder(strings.NewReader(s)).Decode(&schema); err != nil {
		return "", err
	}

	return schema.AccessToken, nil
}
