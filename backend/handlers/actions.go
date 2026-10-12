package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/cg3rard/AoLDatabaseTechnology/backend/middleware"
	"github.com/cg3rard/AoLDatabaseTechnology/backend/models"
)

// POST /api/invoices/:id/payment-proof  (penghuni mengunggah bukti bayar untuk invoice unitnya)
func SubmitPaymentProof(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cl := middleware.GetClaims(c)
		id, err := c.ParamsInt("id")
		if err != nil || id <= 0 {
			return errJSON(c, 400, "id harus berupa angka positif")
		}
		var in struct {
			PaymentProofURL string `json:"payment_proof_url"`
		}
		if err := c.BodyParser(&in); err != nil {
			return errJSON(c, 400, "format JSON tidak valid")
		}
		url := strings.TrimSpace(in.PaymentProofURL)
		if len(url) == 0 || len(url) > 2000 ||
			!(strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://")) {
			return errJSON(c, 400, "payment_proof_url harus berupa URL http(s)")
		}

		var inv models.Invoice
		if err := db.Where("id = ? AND unit_id = ?", id, cl.UnitID).First(&inv).Error; err != nil {
			return writeError(c, err, "gagal mengambil invoice")
		}
		if inv.PaymentStatus == "paid" || inv.PaymentStatus == "cancelled" {
			return errJSON(c, 409, "invoice ini sudah "+map[string]string{"paid": "dibayar", "cancelled": "dibatalkan"}[inv.PaymentStatus])
		}
		err = db.Model(&inv).Updates(map[string]any{
			"payment_proof_url": url,
			"payment_status":    "pending_verification",
		}).Error
		if err != nil {
			return writeError(c, err, "gagal menyimpan bukti bayar")
		}
		return c.JSON(inv)
	}
}

// PATCH /api/facility-bookings/:id/cancel  (pemilik booking atau admin membatalkan booking)
func CancelBooking(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cl := middleware.GetClaims(c)
		id, err := c.ParamsInt("id")
		if err != nil || id <= 0 {
			return errJSON(c, 400, "id harus berupa angka positif")
		}
		q := db.Where("id = ?", id)
		if cl.Role != middleware.RoleAdmin {
			q = q.Where("user_id = ?", cl.UserID)
		}
		var b models.FacilityBooking
		if err := q.First(&b).Error; err != nil {
			return writeError(c, err, "gagal mengambil booking")
		}
		if b.Status != "confirmed" {
			return errJSON(c, 409, "hanya booking berstatus confirmed yang bisa dibatalkan")
		}
		if err := db.Model(&b).Update("status", "cancelled").Error; err != nil {
			return writeError(c, err, "gagal membatalkan booking")
		}
		return c.JSON(b)
	}
}
