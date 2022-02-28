package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Start(errorCh chan<- error) {
	app := fiber.New()

	app.Get("/api/target", listTargets)
	app.Post("/api/target", createTarget)
	app.Put("/api/target", updateTarget)
	app.Get("/api/targetsWithPendingBackups", pendingBackups)

	app.Post("/api/backup", createBackup)
	app.Get("/api/backup/:name", listBackups)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	errorCh <- app.Listen(":3000")
}
