package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/mthuberty/movie-spots-api/internal/db"
	"github.com/mthuberty/movie-spots-api/internal/errs"
	"github.com/mthuberty/movie-spots-api/internal/restclient"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type AppConfiguration struct {
	DBClient   db.DBIface
	OAuthToken string
	RestClient restclient.RestClientIface
	URLS       struct{ VelmatSearchAPI string }
}

var AppConfig *AppConfiguration

func init() {

	AppConfig = &AppConfiguration{}

	isLocal := os.Getenv("GO_ENV") != "nonprod" && os.Getenv("GO_ENV") != "prod"

	if isLocal {
		log.Infoln("Loading local configuration...")

		// Workaround to be able to find .env during tests. https://github.com/joho/godotenv/issues/43
		re := regexp.MustCompile(`^(.*movie-spots-api)`)
		cwd, _ := os.Getwd()
		rootPath := re.Find([]byte(cwd))

		err := godotenv.Load(string(rootPath) + `/.env`)

		if err != nil {
			log.Fatal(errs.WrapTrace("config", "init", fmt.Errorf("Error loading .env file - %s", err)))
		}
	}

	port := os.Getenv("POSTGRES_PORT")
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DATABASE")

	dbc, err := initDB(port, host, user, password, database)
	if err != nil {
		log.Fatal(errs.WrapTrace("config", "init", err))
	}

	AppConfig.DBClient = dbc

	if AppConfig.DBClient == nil {
		log.Fatal(errs.WrapTrace("config", "init", errors.New("Failed to initialize client, DBClient is nil")))
	}
}

func initDB(port, host, user, password, database string) (db.DBIface, error) {

	dbc, err := db.NewDB(port, host, user, password, database)
	if err != nil {
		return nil, errs.WrapTrace("config", "initDB", err)
	}

	return dbc, nil
}
