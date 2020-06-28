package handler

import (
	"net/http"

	"github.com/mthuberty/movie-spots-api/internal/errs"
)

type Handler struct {
	HandlerFunc func(w http.ResponseWriter, r *http.Request) error
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.HandlerFunc(w, r)
	if err != nil {
		errs.HandleErrorResponse(w, err)
	}
}
