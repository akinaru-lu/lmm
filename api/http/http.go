package http

import (
	"fmt"
	"log"
	"net/http"
)

const (
	StatusOK                  = http.StatusOK
	StatusCreated             = http.StatusCreated
	StatusNoContent           = http.StatusNoContent
	StatusBadRequest          = http.StatusBadRequest
	StatusUnauthorized        = http.StatusUnauthorized
	StatusForbidden           = http.StatusForbidden
	StatusNotFound            = http.StatusNotFound
	StatusInternalServerError = http.StatusInternalServerError
)

type Handler = func(Context)

func Serve(addr string, r *Router) {
	fmt.Println("Serving at: " + addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func HandleStatus(c Context, code int) {
	c.String(code, StatusText(code))
}

func NoContent(c Context) {
	HandleStatus(c, http.StatusNoContent)
}

func BadRequest(c Context) {
	HandleStatus(c, StatusBadRequest)
}

func Unauthorized(c Context) {
	HandleStatus(c, StatusUnauthorized)
}

func NotFound(c Context) {
	HandleStatus(c, StatusNotFound)
}

func InternalServerError(c Context) {
	HandleStatus(c, StatusInternalServerError)
}

func StatusText(code int) string {
	return http.StatusText(code)
}
