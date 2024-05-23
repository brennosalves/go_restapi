package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {

	// HEARTBEAT
	app.Get("/status", func(c *fiber.Ctx) error {
		currentTime := time.Now()
		return c.JSON(fiber.Map{
			"status":  "OK",
			"message": "Server is running",
			"date":    currentTime.Format("2006-01-02 15:04:05"),
		})
	})
}
