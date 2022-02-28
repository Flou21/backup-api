package main

import (
	"github.com/Coflnet/db-backup/backup-api/api"
	"github.com/Coflnet/db-backup/backup-api/db"
	"github.com/rs/zerolog/log"
)

func main() {

	errorCh := make(chan error)

	log.Info().Msgf("connecting to database..")
	err := db.Connect()
	if err != nil {
		log.Fatal().Err(err).Msgf("error connecting to database")
	}
	defer db.Disconnect()

	log.Info().Msgf("starting server..")
	go api.Start(errorCh)

	err = <-errorCh
	log.Fatal().Err(err).Msgf("fatal error occured")
}
