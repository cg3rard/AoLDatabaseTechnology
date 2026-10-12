package database

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/cg3rard/AoLDatabaseTechnology/backend/models"
)

// SeedAdmin membuat akun admin pertama jika belum ada admin sama sekali.
// Tidak melakukan apa pun jika email/password kosong. Mengembalikan true jika akun dibuat.
func SeedAdmin(db *gorm.DB, email, password string) (bool, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || password == "" {
		return false, nil
	}
	if len(password) < 8 || len(password) > 72 {
		return false, errors.New("ADMIN_PASSWORD harus 8-72 karakter")
	}

	var admins int64
	if err := db.Model(&models.User{}).Where("role = ?", "admin").Count(&admins).Error; err != nil {
		return false, err
	}
	if admins > 0 {
		return false, nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}
	admin := models.User{
		FullName:     "Administrator",
		Email:        email,
		PasswordHash: string(hash),
		PhoneNumber:  "-",
		Role:         "admin",
	}
	if err := db.Create(&admin).Error; err != nil {
		return false, err
	}
	return true, nil
}
