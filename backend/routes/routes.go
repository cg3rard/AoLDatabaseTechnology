package routes

import (
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"gorm.io/gorm"

	"github.com/cg3rard/AoLDatabaseTechnology/backend/handlers"
	"github.com/cg3rard/AoLDatabaseTechnology/backend/middleware"
)

// Singkatan peran supaya daftar akses di bawah mudah dibaca.
const (
	admin    = middleware.RoleAdmin
	resident = middleware.RoleResident
	laundry  = middleware.RoleLaundry
	security = middleware.RoleSecurity
	tech     = middleware.RoleTechnician
)

// register memasang 5 endpoint CRUD untuk satu resource.
//
//	read           peran yang boleh GET (kosong = semua yang sudah login)
//	write          peran yang boleh POST/PUT/DELETE
//	residentCreate true jika penghuni juga boleh POST (data dibatasi ke unitnya lewat OnCreate)
func register[T any](api fiber.Router, authMW fiber.Handler, db *gorm.DB, path string,
	res handlers.Resource[T], read, write []string, residentCreate bool) {

	createRoles := slices.Clone(write)
	if residentCreate {
		createRoles = append(createRoles, resident)
	}

	g := api.Group(path, authMW)
	g.Get("/", middleware.RequireRoles(read...), res.List(db))
	g.Get("/:id", middleware.RequireRoles(read...), res.Get(db))
	g.Post("/", middleware.RequireRoles(createRoles...), res.Create(db))
	g.Put("/:id", middleware.RequireRoles(write...), res.Update(db))
	g.Delete("/:id", middleware.RequireRoles(write...), res.Delete(db))
}

func Setup(app *fiber.App, db *gorm.DB, jwtSecret []byte, tokenTTL time.Duration) {
	api := app.Group("/api")
	authMW := middleware.Auth(jwtSecret)

	// ---- autentikasi ----
	loginLimit := limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "terlalu banyak percobaan login, coba lagi sebentar"})
		},
	})
	api.Post("/auth/login", loginLimit, handlers.Login(db, jwtSecret, tokenTTL))
	api.Get("/auth/me", authMW, handlers.Me(db))

	// ---- pengguna (khusus admin) ----
	users := handlers.UsersResource()
	ug := api.Group("/users", authMW, middleware.RequireRoles(admin))
	ug.Get("/", users.List(db))
	ug.Get("/:id", users.Get(db))
	ug.Post("/", handlers.CreateUser(db))
	ug.Put("/:id", handlers.UpdateUser(db))
	ug.Delete("/:id", handlers.DeleteUser(db))

	// ---- unit (handler tulisan sendiri) ----
	units := api.Group("/units", authMW)
	staff := middleware.RequireRoles(admin, security, laundry, tech)
	units.Get("/", staff, handlers.GetUnits(db))
	units.Get("/:id", staff, handlers.GetUnit(db))
	units.Post("/", middleware.RequireRoles(admin), handlers.CreateUnit(db))
	units.Put("/:id", middleware.RequireRoles(admin), handlers.UpdateUnit(db))
	units.Delete("/:id", middleware.RequireRoles(admin), handlers.DeleteUnit(db))

	all := []string{} // semua peran yang sudah login

	// ---- modul akses ----
	register(api, authMW, db, "/access-cards", handlers.AccessCards(), []string{admin, security, resident}, []string{admin, security}, false)
	register(api, authMW, db, "/card-requests", handlers.CardRequests(), []string{admin, security, resident}, []string{admin, security}, true)

	// ---- modul keuangan ----
	register(api, authMW, db, "/utility-tariffs", handlers.UtilityTariffs(), all, []string{admin}, false)
	register(api, authMW, db, "/invoices", handlers.Invoices(), []string{admin, resident}, []string{admin}, false)
	api.Post("/invoices/:id/payment-proof", authMW, middleware.RequireRoles(resident), handlers.SubmitPaymentProof(db))
	register(api, authMW, db, "/utility-records", handlers.UtilityRecords(), []string{admin, resident}, []string{admin}, false)
	register(api, authMW, db, "/token-purchases", handlers.TokenPurchases(), []string{admin, resident}, []string{admin}, false)

	// ---- modul layanan ----
	register(api, authMW, db, "/packages", handlers.Packages(), []string{admin, security, resident}, []string{admin, security}, false)
	register(api, authMW, db, "/laundry-orders", handlers.LaundryOrders(), []string{admin, laundry, resident}, []string{admin, laundry}, true)
	register(api, authMW, db, "/facilities", handlers.Facilities(), all, []string{admin}, false)
	register(api, authMW, db, "/facility-bookings", handlers.FacilityBookings(), []string{admin, security, resident}, []string{admin}, true)
	api.Patch("/facility-bookings/:id/cancel", authMW, middleware.RequireRoles(admin, resident), handlers.CancelBooking(db))

	// ---- modul operasional ----
	register(api, authMW, db, "/maintenance-tickets", handlers.MaintenanceTickets(), []string{admin, tech, resident}, []string{admin, tech}, true)
	register(api, authMW, db, "/visitor-passes", handlers.VisitorPasses(), []string{admin, security, resident}, []string{admin, security}, true)
	register(api, authMW, db, "/parking-slots", handlers.ParkingSlots(), []string{admin, security, resident}, []string{admin}, false)
	register(api, authMW, db, "/vehicles", handlers.Vehicles(), []string{admin, security, resident}, []string{admin, security}, true)

	// ---- modul komunitas ----
	register(api, authMW, db, "/permits", handlers.Permits(), []string{admin, security, resident}, []string{admin}, true)
	register(api, authMW, db, "/emergency-alerts", handlers.EmergencyAlerts(), []string{admin, security, resident}, []string{admin, security}, true)
	register(api, authMW, db, "/lost-and-found", handlers.LostAndFound(), all, []string{admin, security}, true)
	register(api, authMW, db, "/announcements", handlers.Announcements(), all, []string{admin}, false)
}
