package mapper

import (
	"time"
	api "toll/api/restapi"
	"toll/api/types"
)

func ModelToTollEvent(c *api.TollEvent) *types.TollEvent {
	if c == nil {
		return nil
	}

	ret := &types.TollEvent{
		LicensePlate: c.LicensePlate,
		EventStart:   c.EventStart,
		EventStop:    c.EventStart.Add(time.Hour),
		VehicleType:  types.VehicleType(c.VehicleType),
		Billed:       false,
	}

	return ret
}
