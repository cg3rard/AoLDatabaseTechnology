package models

import "time"

// Catatan umum untuk semua model:
//   - Kolom yang boleh NULL memakai pointer (*string, *uint, *time.Time).
//   - Kolom yang punya DEFAULT di database diberi tag gorm:"default:..." supaya
//     nilai kosong dari klien diganti default database.
//   - BOOLEAN dengan default TRUE memakai *bool, supaya "false" tetap bisa dikirim.

type AccessCard struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	UnitID            uint       `json:"unit_id"`
	RfidUID           string     `json:"rfid_uid" gorm:"column:rfid_uid"`
	CardLabel         string     `json:"card_label"`
	HolderName        string     `json:"holder_name"`
	AuthorizedFloor   int        `json:"authorized_floor"`
	HasFacilityAccess *bool      `json:"has_facility_access" gorm:"default:true"`
	Status            string     `json:"status" gorm:"default:active"`
	IssuedAt          time.Time  `json:"issued_at" gorm:"autoCreateTime"`
	BlockedAt         *time.Time `json:"blocked_at"`
	BlockReason       *string    `json:"block_reason"`
}

func (AccessCard) TableName() string { return "access_cards" }

type CardRequest struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UnitID      uint      `json:"unit_id"`
	RequestedBy uint      `json:"requested_by"`
	RequestType string    `json:"request_type"`
	LostCardID  *uint     `json:"lost_card_id"`
	HolderName  string    `json:"holder_name"`
	Reason      *string   `json:"reason"`
	FeeAmount   float64   `json:"fee_amount"`
	Status      string    `json:"status" gorm:"default:pending"`
	ProcessedBy *uint     `json:"processed_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (CardRequest) TableName() string { return "card_requests" }
