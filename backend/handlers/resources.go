package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/cg3rard/AoLDatabaseTechnology/backend/middleware"
	"github.com/cg3rard/AoLDatabaseTechnology/backend/models"
)

// ---------- helper ----------

func isResident(cl *middleware.Claims) bool { return cl.Role == middleware.RoleResident }

// residentUnit mengembalikan unit milik penghuni, atau error 403 jika akunnya belum terhubung ke unit.
func residentUnit(cl *middleware.Claims) (uint, error) {
	if cl.UnitID == 0 {
		return 0, fiber.NewError(fiber.StatusForbidden, "akun Anda belum terhubung ke unit")
	}
	return cl.UnitID, nil
}

// scopeUnit: penghuni hanya melihat baris milik unitnya (kolom unit_id).
func scopeUnit(cl *middleware.Claims) (string, []any) {
	if isResident(cl) {
		return "unit_id = ?", []any{cl.UnitID}
	}
	return "", nil
}

// scopeOwner: penghuni hanya melihat baris yang kolom `col`-nya berisi id dirinya.
func scopeOwner(col string) func(*middleware.Claims) (string, []any) {
	return func(cl *middleware.Claims) (string, []any) {
		if isResident(cl) {
			return "`" + col + "` = ?", []any{cl.UserID}
		}
		return "", nil
	}
}

// setNow mengisi field waktu (pointer) dengan waktu sekarang jika masih kosong.
func setNow(f **time.Time) {
	if *f == nil {
		t := time.Now()
		*f = &t
	}
}

func randomHex(nBytes int) string {
	b := make([]byte, nBytes)
	_, _ = rand.Read(b)
	return strings.ToUpper(hex.EncodeToString(b))
}

func badRequest(msg string) error { return fiber.NewError(fiber.StatusBadRequest, msg) }

// parseClock membaca "HH:MM" atau "HH:MM:SS" menjadi menit sejak tengah malam.
func parseClock(s string) (int, error) {
	for _, layout := range []string{"15:04:05", "15:04"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.Hour()*60 + t.Minute(), nil
		}
	}
	return 0, errors.New("format jam salah")
}

func formatClock(minutes int) string { return fmt.Sprintf("%02d:%02d:00", minutes/60, minutes%60) }

// ---------- modul akses ----------

func AccessCards() Resource[models.AccessCard] {
	return Resource[models.AccessCard]{
		Filters: []string{"unit_id", "status"},
		Scope:   scopeUnit,
		OnUpdate: func(tx *gorm.DB, cl *middleware.Claims, old, it *models.AccessCard) error {
			if old.Status == "active" && it.Status != "active" {
				setNow(&it.BlockedAt)
			}
			return nil
		},
	}
}

func CardRequests() Resource[models.CardRequest] {
	return Resource[models.CardRequest]{
		Filters: []string{"unit_id", "status", "request_type"},
		Scope:   scopeUnit,
		OnCreate: func(tx *gorm.DB, cl *middleware.Claims, it *models.CardRequest) error {
			if isResident(cl) {
				unit, err := residentUnit(cl)
				if err != nil {
					return err
				}
				it.UnitID = unit
				it.RequestedBy = cl.UserID
				it.Status = "pending"
				it.FeeAmount = 0
				it.ProcessedBy = nil
			} else if it.RequestedBy == 0 {
				it.RequestedBy = cl.UserID
			}
			if it.LostCardID != nil {
				var n int64
				tx.Model(&models.AccessCard{}).Where("id = ? AND unit_id = ?", *it.LostCardID, it.UnitID).Count(&n)
				if n == 0 {
					return badRequest("kartu yang hilang tidak ditemukan di unit ini")
				}
			}
			return nil
		},
	}
}

// ---------- modul keuangan ----------

func UtilityTariffs() Resource[models.UtilityTariff] {
	return Resource[models.UtilityTariff]{Filters: []string{"utility_type"}}
}

func Invoices() Resource[models.Invoice] {
	return Resource[models.Invoice]{
		Filters: []string{"unit_id", "invoice_type", "payment_status"},
		Scope:   scopeUnit,
		OnCreate: func(tx *gorm.DB, cl *middleware.Claims, it *models.Invoice) error {
			if it.InvoiceNumber == "" {
				it.InvoiceNumber = "INV-" + time.Now().Format("20060102") + "-" + randomHex(3)
			}
			return nil
		},
		OnUpdate: func(tx *gorm.DB, cl *middleware.Claims, old, it *models.Invoice) error {
			if old.PaymentStatus != "paid" && it.PaymentStatus == "paid" {
				setNow(&it.PaidAt)
			}
			return nil
		},
	}
}

func UtilityRecords() Resource[models.UtilityRecord] {
	return Resource[models.UtilityRecord]{
		Filters: []string{"unit_id", "utility_type"},
		Scope:   scopeUnit,
		OnCreate: func(tx *gorm.DB, cl *middleware.Claims, it *models.UtilityRecord) error {
			it.InputtedBy = cl.UserID
			if it.FinalMeter < it.InitialMeter {
				return badRequest("final_meter tidak boleh lebih kecil dari initial_meter")
			}
			return nil
		},
		OnUpdate: func(tx *gorm.DB, cl *middleware.Claims, old, it *models.UtilityRecord) error {
			if it.FinalMeter < it.InitialMeter {
				return badRequest("final_meter tidak boleh lebih kecil dari initial_meter")
			}
			return nil
		},
	}
}

func TokenPurchases() Resource[models.TokenPurchase] {
	return Resource[models.TokenPurchase]{
		Filters: []string{"unit_id", "status"},
		Scope:   scopeUnit,
	}
}

// ---------- modul layanan ----------

func Packages() Resource[models.Package] {
	return Resource[models.Package]{
		Filters: []string{"unit_id", "status"},
		Scope:   scopeUnit,
		OnCreate: func(tx *gorm.DB, cl *middleware.Claims, it *models.Package) error {
			it.ReceivedByStaffID = cl.UserID
			return nil
		},
		OnUpdate: func(tx *gorm.DB, cl *middleware.Claims, old, it *models.Package) error {
			if old.Status != "collected" && it.Status == "collected" {
				setNow(&it.CollectedAt)
			}
			return nil
		},
	}
}

func LaundryOrders() Resource[models.LaundryOrder] {
	return Resource[models.LaundryOrder]{
		Filters: []string{"unit_id", "status", "staff_id"},
		Scope:   scopeUnit,
		OnCreate: func(tx *gorm.DB, cl *middleware.Claims, it *models.LaundryOrder) error {
			if isResident(cl) {
				unit, err := residentUnit(cl)
				if err != nil {
					return err
				}
				it.UnitID = unit
				it.Status = "order_placed"
				it.StaffID = nil
				it.InvoiceID = nil
				it.TotalCost = 0
				it.CompletedAt = nil
			}
			return nil
		},
		OnUpdate: func(tx *gorm.DB, cl *middleware.Claims, old, it *models.LaundryOrder) error {
			if cl.Role == middleware.RoleLaundry && it.StaffID == nil {
				id := cl.UserID
				it.StaffID = &id
			}
			if old.Status != "completed" && it.Status == "completed" {
				setNow(&it.CompletedAt)
			}
			return nil
		},
	}
}

func Facilities() Resource[models.Facility] {
	return Resource[models.Facility]{Filters: []string{"is_active"}}
}

func FacilityBookings() Resource[models.FacilityBooking] {
	return Resource[models.FacilityBooking]{
		Filters: []string{"facility_id", "user_id", "booking_date", "status"},
		Scope:   scopeOwner("user_id"),
		OnCreate: func(tx *gorm.DB, cl *middleware.Claims, it *models.FacilityBooking) error {
			if isResident(cl) {
				it.UserID = cl.UserID
				it.Status = "confirmed"
			} else if it.UserID == 0 {
				it.UserID = cl.UserID
			}
			if it.BookingDate.IsZero() {
				return badRequest("booking_date wajib diisi (YYYY-MM-DD)")
			}
			now := time.Now()
			today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
			if isResident(cl) && it.BookingDate.Before(today) {
				return badRequest("tanggal booking tidak boleh di masa lalu")
			}

			// Kunci baris fasilitas: dua booking bersamaan untuk fasilitas yang sama diproses bergantian,
			// sehingga pengecekan kapasitas di bawah tidak bisa "kebobolan".
			var fac models.Facility
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&fac, it.FacilityID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return badRequest("fasilitas tidak ditemukan")
				}
				return err
			}
			if fac.IsActive != nil && !*fac.IsActive {
				return fiber.NewError(fiber.StatusConflict, "fasilitas sedang tidak aktif")
			}

			start, err1 := parseClock(it.StartTime)
			end, err2 := parseClock(it.EndTime)
			open, err3 := parseClock(fac.OpenTime)
			closeAt, err4 := parseClock(fac.CloseTime)
			if err1 != nil || err2 != nil {
				return badRequest("start_time dan end_time harus berformat HH:MM")
			}
			if err3 != nil || err4 != nil {
				return errors.New("jam operasional fasilitas tidak valid")
			}
			if end <= start {
				return badRequest("end_time harus lebih besar dari start_time")
			}
			if start < open || end > closeAt {
				return badRequest("di luar jam operasional fasilitas")
			}
			if end-start != fac.SlotDurationMinutes {
				return badRequest(fmt.Sprintf("durasi booking harus %d menit", fac.SlotDurationMinutes))
			}
			it.StartTime, it.EndTime = formatClock(start), formatClock(end)

			var used int64
			err := tx.Model(&models.FacilityBooking{}).
				Where("facility_id = ? AND booking_date = ? AND status = 'confirmed' AND start_time < ? AND end_time > ?",
					it.FacilityID, it.BookingDate, it.EndTime, it.StartTime).
				Count(&used).Error
			if err != nil {
				return err
			}
			if int(used) >= fac.MaxCapacity {
				return fiber.NewError(fiber.StatusConflict, "slot waktu ini sudah penuh")
			}
			return nil
		},
	}
}

// ---------- modul operasional ----------

func MaintenanceTickets() Resource[models.MaintenanceTicket] {
	return Resource[models.MaintenanceTicket]{
		Filters: []string{"unit_id", "status", "category", "priority", "assigned_technician_id"},
		Scope:   scopeUnit,
		OnCreate: func(tx *gorm.DB, cl *middleware.Claims, it *models.MaintenanceTicket) error {
			it.ReportedBy = cl.UserID
			if isResident(cl) {
				unit, err := residentUnit(cl)
				if err != nil {
					return err
				}
				it.UnitID = &unit
				it.Priority = "medium"
				it.Status = "open"
				it.AssignedTechnicianID = nil
				it.ResolutionNotes = nil
				it.ResolvedAt = nil
			}
			return nil
		},
		OnUpdate: func(tx *gorm.DB, cl *middleware.Claims, old, it *models.MaintenanceTicket) error {
			if cl.Role == middleware.RoleTechnician && it.AssignedTechnicianID == nil {
				id := cl.UserID
				it.AssignedTechnicianID = &id
			}
			done := func(s string) bool { return s == "resolved" || s == "closed" }
			if !done(old.Status) && done(it.Status) {
				setNow(&it.ResolvedAt)
			}
			return nil
		},
	}
}

func VisitorPasses() Resource[models.VisitorPass] {
	return Resource[models.VisitorPass]{
		Filters: []string{"unit_id", "status", "valid_date"},
		Scope:   scopeUnit,
		OnCreate: func(tx *gorm.DB, cl *middleware.Claims, it *models.VisitorPass) error {
			if isResident(cl) {
				unit, err := residentUnit(cl)
				if err != nil {
					return err
				}
				it.UnitID = unit
				it.Status = "active"
				it.CheckInTime = nil
				it.CheckOutTime = nil
				it.PassCode = "" // selalu dibuat server
			}
			if it.PassCode == "" {
				it.PassCode = randomHex(8)
			}
			return nil
		},
	}
}

func ParkingSlots() Resource[models.ParkingSlot] {
	return Resource[models.ParkingSlot]{Filters: []string{"zone", "vehicle_type", "is_occupied"}}
}

func Vehicles() Resource[models.Vehicle] {
	return Resource[models.Vehicle]{
		Filters: []string{"unit_id", "status", "vehicle_type"},
		Scope:   scopeUnit,
		OnCreate: func(tx *gorm.DB, cl *middleware.Claims, it *models.Vehicle) error {
			if isResident(cl) {
				unit, err := residentUnit(cl)
				if err != nil {
					return err
				}
				it.UnitID = unit
				it.ParkingSlotID = nil
				it.Status = "active"
			}
			return nil
		},
	}
}

// ---------- modul komunitas ----------

func Permits() Resource[models.MovingRenovationPermit] {
	check := func(it *models.MovingRenovationPermit) error {
		if it.StartDate.IsZero() || it.EndDate.IsZero() {
			return badRequest("start_date dan end_date wajib diisi (YYYY-MM-DD)")
		}
		if it.EndDate.Before(it.StartDate.Time) {
			return badRequest("end_date tidak boleh sebelum start_date")
		}
		return nil
	}
	return Resource[models.MovingRenovationPermit]{
		Filters: []string{"unit_id", "status", "permit_type"},
		Scope:   scopeUnit,
		OnCreate: func(tx *gorm.DB, cl *middleware.Claims, it *models.MovingRenovationPermit) error {
			if isResident(cl) {
				unit, err := residentUnit(cl)
				if err != nil {
					return err
				}
				it.UnitID = unit
				it.Status = "submitted"
				it.DepositAmount = 0
				it.DepositInvoiceID = nil
			}
			return check(it)
		},
		OnUpdate: func(tx *gorm.DB, cl *middleware.Claims, old, it *models.MovingRenovationPermit) error {
			return check(it)
		},
	}
}

func EmergencyAlerts() Resource[models.EmergencyAlert] {
	return Resource[models.EmergencyAlert]{
		Filters: []string{"unit_id", "status", "alert_type"},
		Scope:   scopeUnit,
		OnCreate: func(tx *gorm.DB, cl *middleware.Claims, it *models.EmergencyAlert) error {
			it.TriggeredBy = cl.UserID
			if isResident(cl) {
				unit, err := residentUnit(cl)
				if err != nil {
					return err
				}
				it.UnitID = unit
				it.Status = "triggered"
				it.HandledBy = nil
				it.ResolvedAt = nil
			}
			return nil
		},
		OnUpdate: func(tx *gorm.DB, cl *middleware.Claims, old, it *models.EmergencyAlert) error {
			if it.HandledBy == nil && it.Status != "triggered" {
				id := cl.UserID
				it.HandledBy = &id
			}
			done := func(s string) bool { return s == "resolved" || s == "false_alarm" }
			if !done(old.Status) && done(it.Status) {
				setNow(&it.ResolvedAt)
			}
			return nil
		},
	}
}

func LostAndFound() Resource[models.LostAndFound] {
	return Resource[models.LostAndFound]{
		Filters: []string{"type", "status", "category"},
		OnCreate: func(tx *gorm.DB, cl *middleware.Claims, it *models.LostAndFound) error {
			it.ReportedBy = cl.UserID
			if isResident(cl) {
				it.Status = "open"
				it.ClaimedBy = nil
				it.ClaimedAt = nil
			}
			return nil
		},
		OnUpdate: func(tx *gorm.DB, cl *middleware.Claims, old, it *models.LostAndFound) error {
			if old.Status != "claimed" && it.Status == "claimed" {
				setNow(&it.ClaimedAt)
			}
			return nil
		},
	}
}

func Announcements() Resource[models.Announcement] {
	return Resource[models.Announcement]{
		Filters: []string{"category", "is_published"},
		Order:   "id DESC",
		// Selain admin hanya melihat pengumuman yang sudah dipublikasikan.
		Scope: func(cl *middleware.Claims) (string, []any) {
			if cl.Role != middleware.RoleAdmin {
				return "is_published = ?", []any{true}
			}
			return "", nil
		},
		OnCreate: func(tx *gorm.DB, cl *middleware.Claims, it *models.Announcement) error {
			it.AuthorID = cl.UserID
			return nil
		},
	}
}
