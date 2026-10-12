package models

import "time"

type Package struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	UnitID            uint       `json:"unit_id"`
	RecipientName     string     `json:"recipient_name"`
	CourierName       string     `json:"courier_name"`
	ReceiptNumber     *string    `json:"receipt_number"`
	PackagePhotoURL   string     `json:"package_photo_url"`
	Status            string     `json:"status" gorm:"default:waiting_pickup"`
	ReceivedByStaffID uint       `json:"received_by_staff_id"`
	CollectedAt       *time.Time `json:"collected_at"`
	ReceivedAt        time.Time  `json:"received_at" gorm:"autoCreateTime"`
}

func (Package) TableName() string { return "packages" }

type LaundryOrder struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	UnitID      uint       `json:"unit_id"`
	StaffID     *uint      `json:"staff_id"`
	InvoiceID   *uint      `json:"invoice_id"`
	ServiceType string     `json:"service_type"`
	WeightKg    *float64   `json:"weight_kg"`
	ItemCount   *int       `json:"item_count"`
	TotalCost   float64    `json:"total_cost"`
	Status      string     `json:"status" gorm:"default:order_placed"`
	Notes       *string    `json:"notes"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

func (LaundryOrder) TableName() string { return "laundry_orders" }

// OpenTime/CloseTime/StartTime/EndTime bertipe string karena kolomnya TIME ("HH:MM:SS").
type Facility struct {
	ID                  uint   `json:"id" gorm:"primaryKey"`
	Name                string `json:"name"`
	MaxCapacity         int    `json:"max_capacity"`
	OpenTime            string `json:"open_time"`
	CloseTime           string `json:"close_time"`
	SlotDurationMinutes int    `json:"slot_duration_minutes"`
	IsActive            *bool  `json:"is_active" gorm:"default:true"`
}

func (Facility) TableName() string { return "facilities" }

type FacilityBooking struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FacilityID  uint      `json:"facility_id"`
	UserID      uint      `json:"user_id"`
	BookingDate Date      `json:"booking_date"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	Status      string    `json:"status" gorm:"default:confirmed"`
	CreatedAt   time.Time `json:"created_at"`
}

func (FacilityBooking) TableName() string { return "facility_bookings" }
