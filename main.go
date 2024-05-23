package main

import (
	"log"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/brennosalves/go_restapi/config"
	"github.com/brennosalves/go_restapi/services/router"
)

func main() {

	// LOAD API CONFIGURATION
	if resp := config.Init(); resp != nil {
		log.Fatalf("Erro na leitura das configurações da API:%v", resp)
	}

	// NUMBER OF CORES
	runtime.GOMAXPROCS(config.ApiCores)

	// INICIALIZAÇÃO DA API
	app := fiber.New(fiber.Config{
		Prefork:               true,
		AppName:               config.ApiName,
		DisableStartupMessage: false,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET",
	}))

	// FIBER LOGGER FORMAT
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${method} | ${status} | ${path} | ${ip} | ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05.000",
	}))

	// SETUP ROUTER
	router.SetupRouter(app)

	// START
	log.Fatalf("Erro ao iniciar API: %v", app.Listen(":"+config.ApiPort))
}
