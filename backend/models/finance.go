package models

import "time"

// Nominal uang memakai float64 agar mudah dipelajari.
// Untuk sistem keuangan sungguhan, pakai tipe decimal.

type UtilityTariff struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UtilityType string    `json:"utility_type"`
	RatePerUnit float64   `json:"rate_per_unit"`
	AdminFee    float64   `json:"admin_fee"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (UtilityTariff) TableName() string { return "utility_tariffs" }

type Invoice struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	InvoiceNumber   string     `json:"invoice_number"`
	UnitID          uint       `json:"unit_id"`
	InvoiceType     string     `json:"invoice_type"`
	TotalAmount     float64    `json:"total_amount"`
	PaymentStatus   string     `json:"payment_status" gorm:"default:unpaid"`
	PaymentProofURL *string    `json:"payment_proof_url"`
	DueDate         Date       `json:"due_date"`
	PaidAt          *time.Time `json:"paid_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

func (Invoice) TableName() string { return "invoices" }

type UtilityRecord struct {
	ID           uint    `json:"id" gorm:"primaryKey"`
	UnitID       uint    `json:"unit_id"`
	UtilityType  string  `json:"utility_type"`
	InitialMeter float64 `json:"initial_meter"`
	FinalMeter   float64 `json:"final_meter"`
	// usage_total dihitung otomatis oleh database (GENERATED), jadi hanya-baca.
	UsageTotal  float64 `json:"usage_total" gorm:"<-:false"`
	ReadingDate Date    `json:"reading_date"`
	InvoiceID   *uint   `json:"invoice_id"`
	InputtedBy  uint    `json:"inputted_by"`
}

func (UtilityRecord) TableName() string { return "utility_records" }

type TokenPurchase struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UnitID    uint      `json:"unit_id"`
	InvoiceID uint      `json:"invoice_id"`
	Nominal   float64   `json:"nominal"`
	TokenCode *string   `json:"token_code"`
	Status    string    `json:"status" gorm:"default:pending"`
	CreatedAt time.Time `json:"created_at"`
}

func (TokenPurchase) TableName() string { return "token_purchases" }
