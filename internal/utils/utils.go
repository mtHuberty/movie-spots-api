package utils

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

func PrettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Error(err)
		return
	}
	log.Infoln(string(b))
}
