package types

import (
	"time"
)

// VehicleType is a string-based enum.
type VehicleType string

// Enum values
const (
	Motorbike VehicleType = "motorbike"
	Tractor   VehicleType = "tractor"
	Emergency VehicleType = "emergency"
	Diplomat  VehicleType = "diplomat"
	Foreign   VehicleType = "foreign"
	Military  VehicleType = "military"
)

type TollEvent struct {
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	LicensePlate string      `json:"license_plate" db:"license_plate"`
	EventStart   time.Time   `json:"event_start" db:"event_start"`
	EventStop    time.Time   `json:"event_stop" db:"event_stop"`
	VehicleType  VehicleType `json:"vehicle_type" db:"vehicle_type"`
	Billed       bool        `json:"billed" db:"billed"`
}

// IsTollFree checks if a vehicle is toll-free.
func (t TollEvent) IsTollFree() bool {
	switch t.VehicleType {
	case Motorbike, Tractor, Emergency, Diplomat, Foreign, Military:
		return true
	default:
		return false
	}
}

// DailyFee holds total daily amount for specific license plate.
type DailyFee struct {
	Date         time.Time `json:"date" db:"date"`
	LicensePlate string    `json:"license_plate" db:"license_plate"`
	Fee          int       `json:"fee" db:"fee"`
}
