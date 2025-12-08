package service

import (
	"testing"
	"time"
	"toll/api/types"
	"toll/internal/log"
	"toll/internal/test"
)

func TestDeviceDataServiceGetConnectedDevices(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		events      []*types.TollEvent
		expectedFee int
	}{
		{
			"simple event case",
			[]*types.TollEvent{
				{
					EventStart:  time.Date(2025, time.February, 3, 7, 0, 0, 0, time.UTC), // 7:00 - 18
					VehicleType: types.Car,
				},
			},
			18,
		},
		{
			"toll free vehicle",
			[]*types.TollEvent{
				{
					EventStart:  time.Date(2025, time.February, 3, 7, 0, 0, 0, time.UTC),
					VehicleType: types.Tractor,
				},
			},
			0,
		},
		{
			"toll free event day",
			[]*types.TollEvent{
				{
					EventStart:  time.Date(2025, time.January, 1, 8, 0, 0, 0, time.UTC), // New Year's Day
					VehicleType: types.Car,
				},
			},
			0,
		},
		{
			"multiple events in same 60-min window (take highest)",
			[]*types.TollEvent{
				{EventStart: time.Date(2025, 2, 3, 7, 0, 0, 0, time.UTC), VehicleType: types.Car},  // 18
				{EventStart: time.Date(2025, 2, 3, 7, 20, 0, 0, time.UTC), VehicleType: types.Car}, // 18 - ignored
				{EventStart: time.Date(2025, 2, 3, 7, 50, 0, 0, time.UTC), VehicleType: types.Car}, // 18 - ignored
			},
			18, // only max in window
		},
		{
			"same window fee replacement",
			[]*types.TollEvent{
				{EventStart: time.Date(2025, 2, 3, 6, 5, 0, 0, time.UTC), VehicleType: types.Car},  // 6:05 - 8
				{EventStart: time.Date(2025, 2, 3, 6, 35, 0, 0, time.UTC), VehicleType: types.Car}, // 6:35 - 13, replaces 8
			},
			13, // only the max in the window counts
		},
		{
			"events in separate windows",
			[]*types.TollEvent{
				{EventStart: time.Date(2025, 2, 3, 6, 45, 0, 0, time.UTC), VehicleType: types.Car}, // 6:45 - 13
				{EventStart: time.Date(2025, 2, 3, 8, 0, 0, 0, time.UTC), VehicleType: types.Car},  // 8:00 - 13
			},
			26,
		},
		{
			"mixed free & paid – free event ignored",
			[]*types.TollEvent{
				{EventStart: time.Date(2025, 2, 3, 7, 0, 0, 0, time.UTC), VehicleType: types.Car},      // 18
				{EventStart: time.Date(2025, 2, 3, 7, 30, 0, 0, time.UTC), VehicleType: types.Tractor}, // free
			},
			18,
		},
		{
			"no events",
			[]*types.TollEvent{},
			0,
		},
		{
			"daily cap reached",
			[]*types.TollEvent{
				{EventStart: time.Date(2025, 2, 3, 6, 0, 0, 0, time.UTC), VehicleType: types.Car},  // 8
				{EventStart: time.Date(2025, 2, 3, 7, 31, 0, 0, time.UTC), VehicleType: types.Car}, // 7:31 - 18
				{EventStart: time.Date(2025, 2, 3, 9, 1, 0, 0, time.UTC), VehicleType: types.Car},  // 9:01 - 8
				{EventStart: time.Date(2025, 2, 3, 11, 0, 0, 0, time.UTC), VehicleType: types.Car}, // 8
				{EventStart: time.Date(2025, 2, 3, 13, 0, 0, 0, time.UTC), VehicleType: types.Car}, // 8
				{EventStart: time.Date(2025, 2, 3, 15, 0, 0, 0, time.UTC), VehicleType: types.Car}, // 15:00 - 13
				{EventStart: time.Date(2025, 2, 3, 17, 0, 0, 0, time.UTC), VehicleType: types.Car}, // 17:00 - 13
				{EventStart: time.Date(2025, 2, 3, 18, 1, 0, 0, time.UTC), VehicleType: types.Car}, // 18:01 - 0
			},
			maxDailyFee,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := &billing{
				log: log.Noop(),
			}

			totalFeeResult := svc.calculateDailyFee(tc.events)

			test.Match(t, totalFeeResult, tc.expectedFee)
		})
	}
}
