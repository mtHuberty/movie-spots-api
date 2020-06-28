package config

import (
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

		// Find .env during tests. Workaround from https://github.com/joho/godotenv/issues/43
		re := regexp.MustCompile(`^(.*movie-spots-api)`)
		cwd, _ := os.Getwd()
		rootPath := re.Find([]byte(cwd))

		err := godotenv.Load(string(rootPath) + `/.env`)

		if err != nil {
			log.Fatal(errs.WrapTrace("config", "init", fmt.Errorf("Error loading .env file - %s", err)))
		}
	}
}

func initDB(postgresSecrets map[string]string, isLocal bool) (db.DBIface, error) {
	portstring := postgresSecrets["port"]
	// override the pg port if we are running locally and using the tunnel
	// to connect to the db.
	if pgPort := os.Getenv("POSTGRES_PORT"); pgPort != "" && isLocal {
		log.Infof("Using local config for pg_port: %s \n", pgPort)
		portstring = pgPort
	}

	dbc, err := db.NewDB(postgresSecrets["host"], portstring, postgresSecrets["user"], postgresSecrets["password"], postgresSecrets["database"])
	if err != nil {
		return nil, errs.WrapTrace("config", "initDB", err)
	}

	return dbc, nil
}
