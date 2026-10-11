package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cg3rard/AoLDatabaseTechnology/backend/database"
	"github.com/cg3rard/AoLDatabaseTechnology/backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Environment variable %s belum diisi (lihat .env.example)", key)
	}
	return v
}

func main() {
	_ = godotenv.Load("../.env", ".env")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mustEnv("DB_USER"),
		mustEnv("DB_PASSWORD"),
		getEnv("DB_HOST", "127.0.0.1"),
		getEnv("DB_PORT", "3306"),
		mustEnv("DB_NAME"),
	)
	db, err := database.Connect(dsn)
	if err != nil {
		log.Fatal("Gagal terhubung ke database: ", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Gagal mengambil koneksi database: ", err)
	}
	defer sqlDB.Close()
	log.Println("Berhasil terhubung ke MySQL!")

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: getEnv("CORS_ORIGINS", "http://localhost:5173"),
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Backend API Property Management Berjalan!"})
	})

	app.Get("/api/health", func(c *fiber.Ctx) error {
		if err := sqlDB.Ping(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "db_down"})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})

	routes.Setup(app, db)

	port := getEnv("APP_PORT", "8080")
	log.Fatal(app.Listen("127.0.0.1:" + port))
}