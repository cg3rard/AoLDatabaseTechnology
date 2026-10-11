package models

import "time"

// Unit mewakili satu baris di tabel `units`.
// GORM membaca nama struct "Unit" sebagai tabel "units",
// dan nama field "UnitNumber" sebagai kolom "unit_number".
type Unit struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UnitNumber string    `json:"unit_number"`
	TowerBlock string    `json:"tower_block"`
	Floor      int       `json:"floor"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
