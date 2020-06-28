package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/mthuberty/movie-spots-api/internal/config"
	"github.com/mthuberty/movie-spots-api/internal/router"

	log "github.com/sirupsen/logrus"
)

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3001"
	}

	ac := config.AppConfig
	defer ac.DBClient.Close()

	r := router.NewRouter(ac)

	http.Handle("/", r)

	log.Infof("Listening on port " + PORT + "...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), r))
}
