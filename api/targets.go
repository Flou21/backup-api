package api

import (
	"encoding/json"

	"github.com/Coflnet/db-backup/server/db"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func listTargets(c *fiber.Ctx) error {
	targets, err := db.ListTargets()

	if err != nil {
		log.Error().Err(err).Msgf("listing targets from the database failed")
		return err
	}

	serialized, err := json.Marshal(targets)
	if err != nil {
		log.Error().Err(err).Msgf("failed to marshal targets")
		return err
	}

	return c.Status(fiber.StatusOK).SendString(string(serialized))
}

func pendingBackups(c *fiber.Ctx) error {
	targets := make(chan *db.Target)
	go db.ListTargetsWithPendingBackups(targets)

	slice := []*db.Target{}
	for target := range targets {
		slice = append(slice, target)
	}

	serialized, err := json.Marshal(slice)
	if err != nil {
		log.Error().Err(err).Msgf("failed to marshal targets")
		return err
	}

	return c.Status(fiber.StatusOK).SendString(string(serialized))
}

func createTarget(c *fiber.Ctx) error {

	target := new(db.Target)

	if err := c.BodyParser(target); err != nil {
		log.Error().Err(err).Msgf("invalid body for creating a target: %s", string(c.Body()))
		return err
	}

	log.Info().Msgf("inserting backup target with name %s", target.Name)

	err := target.Insert()
	if err != nil {
		log.Error().Err(err).Msgf("inserting target failed")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func updateTarget(c *fiber.Ctx) error {
	target := new(db.Target)

	if err := c.BodyParser(target); err != nil {
		log.Error().Err(err).Msgf("there was an error when parsing target in put request, body: ", string(c.Body()))
		return err
	}

	err := target.Update()
	if err != nil {
		log.Error().Err(err).Msgf("an error occurred when updating target in put request")
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
