package handlers

import (
	"errors"

	"github.com/cg3rard/AoLDatabaseTechnology/backend/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var validStatuses = map[string]bool{
	"vacant": true, "occupied": true, "under_renovation": true,
}

// validateUnit mengembalikan pesan error, atau "" kalau data valid.
func validateUnit(u *models.Unit) string {
	if u.UnitNumber == "" {
		return "unit_number wajib diisi"
	}
	if u.TowerBlock == "" {
		return "tower_block wajib diisi"
	}
	if u.Status == "" {
		u.Status = "vacant"
	}
	if !validStatuses[u.Status] {
		return "status harus vacant, occupied, atau under_renovation"
	}
	return ""
}

func errJSON(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(fiber.Map{"error": msg})
}

// GET /api/units
func GetUnits(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var units []models.Unit
		if err := db.Order("id").Find(&units).Error; err != nil {
			return errJSON(c, 500, "gagal mengambil data unit")
		}
		return c.JSON(units)
	}
}

// GET /api/units/:id
func GetUnit(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return errJSON(c, 400, "id harus berupa angka")
		}
		var unit models.Unit
		if err := db.First(&unit, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errJSON(c, 404, "unit tidak ditemukan")
			}
			return errJSON(c, 500, "gagal mengambil unit")
		}
		return c.JSON(unit)
	}
}

// POST /api/units
func CreateUnit(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input models.Unit
		if err := c.BodyParser(&input); err != nil {
			return errJSON(c, 400, "format JSON tidak valid")
		}
		if msg := validateUnit(&input); msg != "" {
			return errJSON(c, 400, msg)
		}
		input.ID = 0 // id selalu dibuat database
		if err := db.Create(&input).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return errJSON(c, 409, "unit_number sudah dipakai")
			}
			return errJSON(c, 500, "gagal menyimpan unit")
		}
		return c.Status(201).JSON(input)
	}
}

// PUT /api/units/:id
func UpdateUnit(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return errJSON(c, 400, "id harus berupa angka")
		}
		var unit models.Unit
		if err := db.First(&unit, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errJSON(c, 404, "unit tidak ditemukan")
			}
			return errJSON(c, 500, "gagal mengambil unit")
		}

		var input models.Unit
		if err := c.BodyParser(&input); err != nil {
			return errJSON(c, 400, "format JSON tidak valid")
		}
		if msg := validateUnit(&input); msg != "" {
			return errJSON(c, 400, msg)
		}

		err = db.Model(&unit).Updates(map[string]any{
			"unit_number": input.UnitNumber,
			"tower_block": input.TowerBlock,
			"floor":       input.Floor,
			"status":      input.Status,
		}).Error
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return errJSON(c, 409, "unit_number sudah dipakai")
			}
			return errJSON(c, 500, "gagal mengubah unit")
		}
		db.First(&unit, id) // ambil ulang agar response berisi data terbaru
		return c.JSON(unit)
	}
}

// DELETE /api/units/:id
func DeleteUnit(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return errJSON(c, 400, "id harus berupa angka")
		}
		result := db.Delete(&models.Unit{}, id)
		if result.Error != nil {
			return errJSON(c, 500, "gagal menghapus unit")
		}
		if result.RowsAffected == 0 {
			return errJSON(c, 404, "unit tidak ditemukan")
		}
		return c.SendStatus(204)
	}
}
