package handlers

import (
	"errors"
	"net/mail"
	"slices"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/cg3rard/AoLDatabaseTechnology/backend/middleware"
	"github.com/cg3rard/AoLDatabaseTechnology/backend/models"
)

// dummyHash dipakai agar waktu respons login sama baik email ada maupun tidak
// (mencegah penyerang menebak email yang terdaftar dari lamanya respons).
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("password-palsu-untuk-timing"), bcrypt.DefaultCost)

// POST /api/auth/login
func Login(db *gorm.DB, secret []byte, ttl time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&in); err != nil {
			return errJSON(c, 400, "format JSON tidak valid")
		}
		email := strings.ToLower(strings.TrimSpace(in.Email))
		if email == "" || in.Password == "" {
			return errJSON(c, 400, "email dan password wajib diisi")
		}

		var user models.User
		err := db.Where("email = ?", email).First(&user).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return errJSON(c, 500, "gagal memeriksa akun")
			}
			_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(in.Password))
			return errJSON(c, 401, "email atau password salah")
		}
		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)) != nil {
			return errJSON(c, 401, "email atau password salah")
		}

		var unitID uint
		if user.UnitID != nil {
			unitID = *user.UnitID
		}
		token, err := middleware.IssueToken(secret, ttl, user.ID, user.Role, unitID)
		if err != nil {
			return errJSON(c, 500, "gagal membuat token")
		}
		return c.JSON(fiber.Map{"token": token, "user": user})
	}
}

// GET /api/auth/me
func Me(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cl := middleware.GetClaims(c)
		var user models.User
		if err := db.First(&user, cl.UserID).Error; err != nil {
			return errJSON(c, 401, "akun tidak ditemukan")
		}
		return c.JSON(user)
	}
}

// ---------- manajemen pengguna (khusus admin) ----------

var validRoles = []string{
	middleware.RoleAdmin, middleware.RoleResident, middleware.RoleLaundry,
	middleware.RoleSecurity, middleware.RoleTechnician,
}

type userInput struct {
	UnitID       *uint   `json:"unit_id"`
	FullName     string  `json:"full_name"`
	Email        string  `json:"email"`
	Password     string  `json:"password"`
	PhoneNumber  string  `json:"phone_number"`
	Role         string  `json:"role"`
	ResidentType *string `json:"resident_type"`
}

// validate merapikan input dan mengembalikan pesan error ("" jika valid).
func (in *userInput) validate(passwordRequired bool) string {
	in.FullName = strings.TrimSpace(in.FullName)
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	in.PhoneNumber = strings.TrimSpace(in.PhoneNumber)

	if in.FullName == "" || in.PhoneNumber == "" {
		return "full_name dan phone_number wajib diisi"
	}
	if addr, err := mail.ParseAddress(in.Email); err != nil || addr.Address != in.Email {
		return "email tidak valid"
	}
	if passwordRequired || in.Password != "" {
		if len(in.Password) < 8 {
			return "password minimal 8 karakter"
		}
		if len(in.Password) > 72 { // batas bcrypt
			return "password maksimal 72 karakter"
		}
	}
	if !slices.Contains(validRoles, in.Role) {
		return "role harus salah satu dari: " + strings.Join(validRoles, ", ")
	}
	if in.ResidentType != nil && *in.ResidentType != "owner" && *in.ResidentType != "tenant" {
		return "resident_type harus owner atau tenant"
	}
	if in.Role == middleware.RoleResident && in.UnitID == nil {
		return "penghuni (resident) wajib punya unit_id"
	}
	return ""
}

func UsersResource() Resource[models.User] {
	return Resource[models.User]{Filters: []string{"role", "unit_id"}}
}

// POST /api/users
func CreateUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in userInput
		if err := c.BodyParser(&in); err != nil {
			return errJSON(c, 400, "format JSON tidak valid")
		}
		if msg := in.validate(true); msg != "" {
			return errJSON(c, 400, msg)
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return errJSON(c, 500, "gagal memproses password")
		}
		user := models.User{
			UnitID: in.UnitID, FullName: in.FullName, Email: in.Email,
			PasswordHash: string(hash), PhoneNumber: in.PhoneNumber,
			Role: in.Role, ResidentType: in.ResidentType,
		}
		if err := db.Create(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return errJSON(c, 409, "email sudah terdaftar")
			}
			return writeError(c, err, "gagal menyimpan pengguna")
		}
		return c.Status(201).JSON(user)
	}
}

// PUT /api/users/:id  (password boleh dikosongkan = tidak diganti)
func UpdateUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil || id <= 0 {
			return errJSON(c, 400, "id harus berupa angka positif")
		}
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			return writeError(c, err, "gagal mengambil pengguna")
		}
		var in userInput
		if err := c.BodyParser(&in); err != nil {
			return errJSON(c, 400, "format JSON tidak valid")
		}
		if msg := in.validate(false); msg != "" {
			return errJSON(c, 400, msg)
		}
		if in.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
			if err != nil {
				return errJSON(c, 500, "gagal memproses password")
			}
			user.PasswordHash = string(hash)
		}
		user.UnitID, user.FullName, user.Email = in.UnitID, in.FullName, in.Email
		user.PhoneNumber, user.Role, user.ResidentType = in.PhoneNumber, in.Role, in.ResidentType

		if err := db.Save(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return errJSON(c, 409, "email sudah terdaftar")
			}
			return writeError(c, err, "gagal mengubah pengguna")
		}
		return c.JSON(user)
	}
}

// DELETE /api/users/:id
func DeleteUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil || id <= 0 {
			return errJSON(c, 400, "id harus berupa angka positif")
		}
		if uint(id) == middleware.GetClaims(c).UserID {
			return errJSON(c, 400, "Anda tidak bisa menghapus akun Anda sendiri")
		}
		result := db.Delete(&models.User{}, id)
		if result.Error != nil {
			return writeError(c, result.Error, "gagal menghapus pengguna")
		}
		if result.RowsAffected == 0 {
			return errJSON(c, 404, "pengguna tidak ditemukan")
		}
		return c.SendStatus(204)
	}
}
