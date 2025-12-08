package repository

import (
	"context"
	"time"

	"toll/api/types"

	"toll/internal/database"
	"toll/internal/errlog"
)

type (
	// TollEventRepository interface with method definitions.
	TollEventRepository interface {
		GetAll(ctx context.Context, license string, t time.Time) ([]*types.TollEvent, error)
		Record(ctx context.Context, event *types.TollEvent) error
		UpdateDailyFee(ctx context.Context, dailyFee types.DailyFee) error
	}

	tollEvent struct {
		db database.DB
	}
)

// TollEvent func returns ApiKeyRepository with provided database connection.
func TollEvent(db database.DB) TollEventRepository {
	return &tollEvent{db: db}
}

// GetAll returns all toll events for the given license plate at or after a specific time.
func (r *tollEvent) GetAll(ctx context.Context, license string, t time.Time) ([]*types.TollEvent, error) {
	query := `
		SELECT
			created_at,
			license_plate,
			event_start,
			event_stop,
			vehicle_type,
			billed
		FROM events
		WHERE license_plate = $1
		  AND event_start >= $2
		  AND toll_free != true
		ORDER BY event_start
	`

	var events []*types.TollEvent

	err := r.db.Select(ctx, &events, query, license, t)
	if err != nil {
		return nil, errlog.Error(err)
	}

	return events, nil
}

// Record stores a car toll event.
func (r *tollEvent) Record(ctx context.Context, event *types.TollEvent) error {
	insert := `
		INSERT INTO events (
			created_at,
			license_plate,
			event_start,
			event_stop,
			vehicle_type,
			billed,
			toll_free
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(
		ctx,
		insert,
		time.Now(),
		event.LicensePlate,
		event.EventStart,
		event.EventStop,
		event.VehicleType,
		false,
		event.IsTollFree(),
	)
	if err != nil {
		return errlog.Error(err)
	}

	return nil
}

func (r *tollEvent) UpdateDailyFee(ctx context.Context, dailyFee types.DailyFee) error {
	insert := `
		INSERT INTO daily_toll_fees (date, license_plate, fee)
         	VALUES ($1, $2, $3)
		ON CONFLICT (date, license_plate)
		DO UPDATE SET fee = EXCLUDED.fee
	`

	_, err := r.db.Exec(ctx,
		insert,
		dailyFee.Date,
		dailyFee.LicensePlate,
		dailyFee.Fee,
	)
	return err
}
