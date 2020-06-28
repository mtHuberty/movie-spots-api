package router

import (
	"fmt"
	"net/http"

	"github.com/mthuberty/movie-spots-api/internal/config"
	"github.com/mthuberty/movie-spots-api/internal/handler"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func NewRouter(ac *config.AppConfiguration) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/ping", PingPong)

	r.Handle("/reservations", handler.Handler{HandlerFunc: handler.HandleGetReservations(ac)}).Methods("GET")
	r.Handle("/reservations", handler.Handler{HandlerFunc: handler.HandlePostReservations(ac)}).Methods("POST")

	return r
}

func PingPong(w http.ResponseWriter, r *http.Request) {
	log.Infof("Got a request at %v\n", r.Host)
	fmt.Fprintf(w, "pong")
}
