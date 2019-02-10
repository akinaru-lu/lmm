package middleware

import (
	"lmm/api/http"
	"lmm/api/testing"
	"lmm/api/util/contextutil"
	"lmm/api/util/uuidutil"
)

func TestWithRequestID(tt *testing.T) {
	t := testing.NewTester(tt)

	uuid := uuidutil.New()
	sig := uuidutil.New()

	router := http.NewRouter()
	router.Use(WithRequestID)
	router.GET("/", func(c http.Context) {
		reqID := contextutil.RequestID(c)
		t.Is(uuid, reqID)
		c.String(http.StatusOK, sig)
	})

	tt.Run("WithRequestID", func(tt *testing.T) {
		t := testing.NewTester(tt)
		req := testing.GET("/", nil)
		req.Header.Set("X-Request-ID", uuid)
		res := testing.DoRequest(req, router)

		t.Is(http.StatusOK, res.StatusCode())
		t.Is(sig, res.Body())
	})

	tt.Run("WithoutRequestID", func(tt *testing.T) {
		t := testing.NewTester(tt)
		req := testing.GET("/", nil)
		res := testing.DoRequest(req, router)

		t.Is(http.StatusBadRequest, res.StatusCode())
		t.Is(http.StatusText(http.StatusBadRequest), res.Body())
	})
}
