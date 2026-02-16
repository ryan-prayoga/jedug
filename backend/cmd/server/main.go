package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"jedug-backend/internal/config"
	"jedug-backend/internal/database"
	"jedug-backend/internal/handlers"
	"jedug-backend/internal/repository"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()

	pool := database.TryConnect(cfg.DB)
	defer database.Close()

	app := fiber.New(fiber.Config{
		AppName: "Jedug Backend v1.0",
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, http://localhost:5174, http://localhost:3000",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Jedug Backend is running! 🚀")
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/health", func(c *fiber.Ctx) error {
		dbStatus := "disconnected"
		if pool != nil {
			if err := pool.Ping(c.Context()); err == nil {
				dbStatus = "connected"
			}
		}
		status := "ok"
		if dbStatus != "connected" {
			status = "degraded"
		}
		return c.JSON(fiber.Map{
			"status":  status,
			"message": "Service is running",
			"db":      dbStatus,
		})
	})

	// District routes
	if pool != nil {
		districtRepo := repository.NewDistrictRepo(pool)
		districtHandler := handlers.NewDistrictHandler(districtRepo)
		v1.Get("/districts", districtHandler.GetDistricts)
	}

	port := cfg.Port

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	_ = app.Shutdown()
}
