package routes

import (
	"github.com/cg3rard/AoLDatabaseTechnology/backend/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {
	units := app.Group("/api/units")
	units.Get("/", handlers.GetUnits(db))
	units.Get("/:id", handlers.GetUnit(db))
	units.Post("/", handlers.CreateUnit(db))
	units.Put("/:id", handlers.UpdateUnit(db))
	units.Delete("/:id", handlers.DeleteUnit(db))
}