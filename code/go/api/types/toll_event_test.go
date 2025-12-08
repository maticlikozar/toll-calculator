package types

import (
	"testing"
	"toll/internal/test"
)

func TestVehicleType_IsTollFree(t *testing.T) {
	t.Parallel()

	tests := []struct {
		vehicleType VehicleType
		expected    bool
	}{
		{Car, false},
		{Van, false},
		{Truck, false},
		{Motorbike, true},
		{Tractor, true},
		{Emergency, true},
		{Diplomat, true},
		{Foreign, true},
		{Military, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.vehicleType), func(t *testing.T) {
			t.Parallel()

			val := TollEvent{
				LicensePlate: "test",
				VehicleType:  tt.vehicleType,
			}

			test.Match(t, val.VehicleType, tt.expected, val.IsTollFree())
		})
	}
}
