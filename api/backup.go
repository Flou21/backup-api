package api

import (
	"encoding/json"

	"github.com/Coflnet/db-backup/server/db"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func createBackup(c *fiber.Ctx) error {

	backup := new(db.Backup)

	if err := c.BodyParser(backup); err != nil {
		log.Error().Err(err).Msgf("invalid body for creating a backup: %s", string(c.Body()))
		return err
	}

	log.Info().Msgf("inserting backup with name %s", backup.Target.Name)

	err := backup.Insert()
	if err != nil {
		log.Error().Err(err).Msgf("inserting backup failed")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func listBackups(c *fiber.Ctx) error {

	name := c.Params("name")

	backups, err := db.ListBackups(name)
	if err != nil {
		log.Error().Err(err).Msgf("error when listing backups")
		return err
	}

	serialized, err := json.Marshal(backups)
	if err != nil {
		log.Error().Err(err).Msgf("error when serializing backups")
		return err
	}

	return c.Status(fiber.StatusOK).SendString(string(serialized))
}
