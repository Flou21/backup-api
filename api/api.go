package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Start(errorCh chan<- error) {
	app := fiber.New()

	app.Get("/target", listTargets)
	app.Post("/target", createTarget)
	app.Put("/target", updateTarget)
	app.Get("/targetsWithPendingBackups", pendingBackups)

	app.Post("/backup", createBackup)
	app.Get("/backup/:name", listBackups)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3333",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	errorCh <- app.Listen(":3000")
}
