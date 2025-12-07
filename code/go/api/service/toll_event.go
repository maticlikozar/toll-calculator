package service

import (
	"context"

	"toll/internal/database"

	"toll/api/repository"
	"toll/api/types"
)

type (
	// TollEventService interface with method definitions.
	TollEventService interface {
		Record(ctx context.Context, event *types.TollEvent) error
	}

	event struct {
		events repository.TollEventRepository
	}
)

// TollEvent func returns new TollEventService with new background context.
func TollEvent() TollEventService {
	db := database.Get()

	return &event{
		events: repository.TollEvent(db),
	}
}

// Record func stores all products.
func (svc *event) Record(ctx context.Context, event *types.TollEvent) error {
	return svc.events.Record(ctx, event)
}
