package handler

import (
	"context"
	"time"

	apiErrors "toll/api/handler/errors"
	"toll/api/mapper"

	"toll/api/identity"
	"toll/api/restapi"
)

func (s *apiService) RecordTollEvent(ctx context.Context, params *restapi.TollEvent) (restapi.RecordTollEventRes, error) {
	kid := identity.Get(ctx)
	if kid == nil {
		return nil, apiErrors.ErrAPIUnauthorized
	}

	if !params.EventStart.Before(time.Now()) {
		return nil, apiErrors.ErrAPIBadRequest.
			WithDetails("toll event start date should be in the past")
	}

	tollEvent := mapper.ModelToTollEvent(params)

	err := s.tollEvents.Record(ctx, tollEvent)
	if err != nil {
		s.log.Errore(err)

		return nil, apiErrors.ErrAPIInternal
	}

	s.billing.TriggerFor(tollEvent.LicensePlate)

	return &restapi.RecordTollEventNoContent{}, nil
}
