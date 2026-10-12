package models

import "time"

type MovingRenovationPermit struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	UnitID           uint      `json:"unit_id"`
	PermitType       string    `json:"permit_type"`
	StartDate        Date      `json:"start_date"`
	EndDate          Date      `json:"end_date"`
	ContractorName   *string   `json:"contractor_name"`
	WorkerCount      int       `json:"worker_count"`
	DepositAmount    float64   `json:"deposit_amount"`
	DepositInvoiceID *uint     `json:"deposit_invoice_id"`
	Status           string    `json:"status" gorm:"default:submitted"`
	Notes            *string   `json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
}

func (MovingRenovationPermit) TableName() string { return "moving_renovation_permits" }

type EmergencyAlert struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	UnitID      uint       `json:"unit_id"`
	TriggeredBy uint       `json:"triggered_by"`
	AlertType   string     `json:"alert_type"`
	Status      string     `json:"status" gorm:"default:triggered"`
	HandledBy   *uint      `json:"handled_by"`
	CreatedAt   time.Time  `json:"created_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`
}

func (EmergencyAlert) TableName() string { return "emergency_alerts" }

// LostAndFound butuh TableName manual: GORM akan menebak "lost_and_founds".
type LostAndFound struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	ReportedBy    uint       `json:"reported_by"`
	Type          string     `json:"type"`
	ItemName      string     `json:"item_name"`
	Category      string     `json:"category"`
	LocationFound string     `json:"location_found"`
	PhotoURL      *string    `json:"photo_url"`
	Status        string     `json:"status" gorm:"default:open"`
	ClaimedBy     *uint      `json:"claimed_by"`
	ClaimedAt     *time.Time `json:"claimed_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (LostAndFound) TableName() string { return "lost_and_found" }

type Announcement struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	AuthorID    uint      `json:"author_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Category    string    `json:"category" gorm:"default:general"`
	IsPublished *bool     `json:"is_published" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Announcement) TableName() string { return "announcements" }
