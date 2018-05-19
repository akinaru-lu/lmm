package auth

import (
	accountService "lmm/api/context/account/appservice"
	accountRepository "lmm/api/context/account/domain/repository"
	"net/http"
	"regexp"

	"github.com/akinaru-lu/elesion"
)

var (
	PatternBearerAuthorization = regexp.MustCompile(`^Bearer (.+)$`)
)

func BearerAuth(handler elesion.Handler) elesion.Handler {
	return func(c *elesion.Context) {
		auth := c.Request.Header.Get("Authorization")
		mathed := PatternBearerAuthorization.FindStringSubmatch(auth)
		if len(mathed) != 2 {
			c.Status(http.StatusUnauthorized).String(http.StatusText(http.StatusUnauthorized))
			return
		}
		token := mathed[1]
		user, err := accountService.New(accountRepository.New()).VerifyToken(token)
		if err != nil {
			c.Status(http.StatusUnauthorized).String(http.StatusText(http.StatusUnauthorized))
			return
		}
		c.Fields().Set("user", user)
		handler(c)
	}
}
