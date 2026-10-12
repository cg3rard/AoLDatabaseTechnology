package models

import "time"

// User mewakili tabel `users` (profil + akun login).
// PasswordHash tidak pernah dikirim ke klien (json:"-").
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UnitID       *uint     `json:"unit_id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	PhoneNumber  string    `json:"phone_number"`
	Role         string    `json:"role"`
	ResidentType *string   `json:"resident_type"`
	CreatedAt    time.Time `json:"created_at"`
}
