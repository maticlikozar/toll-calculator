package handler

import (
	"context"

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

	tollEvent := mapper.ModelToTollEvent(params)

	err := s.tollEvents.Record(ctx, tollEvent)
	if err != nil {
		s.log.Errore(err)

		return nil, apiErrors.ErrAPIInternal
	}

	return &restapi.RecordTollEventNoContent{}, nil
}
