package models

import "time"

type MaintenanceTicket struct {
	ID                   uint       `json:"id" gorm:"primaryKey"`
	UnitID               *uint      `json:"unit_id"`
	ReportedBy           uint       `json:"reported_by"`
	AssignedTechnicianID *uint      `json:"assigned_technician_id"`
	Category             string     `json:"category"`
	Title                string     `json:"title"`
	Description          string     `json:"description"`
	PhotoURL             *string    `json:"photo_url"`
	Priority             string     `json:"priority" gorm:"default:medium"`
	Status               string     `json:"status" gorm:"default:open"`
	ResolutionNotes      *string    `json:"resolution_notes"`
	CreatedAt            time.Time  `json:"created_at"`
	ResolvedAt           *time.Time `json:"resolved_at"`
}

func (MaintenanceTicket) TableName() string { return "maintenance_tickets" }

type VisitorPass struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	UnitID        uint       `json:"unit_id"`
	VisitorName   string     `json:"visitor_name"`
	VisitorPhone  *string    `json:"visitor_phone"`
	VehicleNumber *string    `json:"vehicle_number"`
	PassCode      string     `json:"pass_code"`
	ValidDate     Date       `json:"valid_date"`
	CheckInTime   *time.Time `json:"check_in_time"`
	CheckOutTime  *time.Time `json:"check_out_time"`
	Status        string     `json:"status" gorm:"default:active"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (VisitorPass) TableName() string { return "visitor_passes" }

type ParkingSlot struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	SlotCode    string `json:"slot_code"`
	Zone        string `json:"zone"`
	VehicleType string `json:"vehicle_type"`
	IsOccupied  bool   `json:"is_occupied"`
}

func (ParkingSlot) TableName() string { return "parking_slots" }

type Vehicle struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UnitID        uint      `json:"unit_id"`
	ParkingSlotID *uint     `json:"parking_slot_id"`
	PlateNumber   string    `json:"plate_number"`
	VehicleType   string    `json:"vehicle_type"`
	BrandModel    *string   `json:"brand_model"`
	Status        string    `json:"status" gorm:"default:active"`
	CreatedAt     time.Time `json:"created_at"`
}

func (Vehicle) TableName() string { return "vehicles" }
