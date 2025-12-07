package handler

import (
	"context"

	apiErrors "toll/api/handler/errors"

	"toll/api/identity"
	"toll/api/restapi"
)

func (s *apiService) RecordTollEvent(ctx context.Context, params *restapi.TollEvent) (restapi.RecordTollEventRes, error) {
	kid := identity.Get(ctx)
	if kid == nil {
		return nil, apiErrors.ErrAPIUnauthorized
	}

	s.log.Warn("RecordTollEvent not implemented yet!")

	return &restapi.RecordTollEventNoContent{}, nil
}
