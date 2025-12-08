package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"

	apiErrors "toll/api/handler/errors"
	"toll/api/restapi"
	"toll/api/service"
	"toll/api/types"

	log "toll/internal/log"
)

type apiService struct {
	log log.Logger

	tollEvents service.TollEventService
	billing    service.BillingService
}

func NewApiHandlers() *apiService {
	return &apiService{
		log: log.WithField(types.LogComponent, "api/handlers"),

		tollEvents: service.TollEvents,
		billing:    service.Billing,
	}
}

func (s *apiService) Health(ctx context.Context) error {
	return nil
}

func (s *apiService) NewError(ctx context.Context, err error) *restapi.ProblemStatusCode {
	//nolint:all
	switch apiErr := err.(type) {
	case *apiErrors.APIError:
		return &restapi.ProblemStatusCode{
			StatusCode: apiErr.StatusCode,
			Response: restapi.Problem{
				Status:    restapi.NewOptInt32(int32(apiErr.StatusCode)),
				Title:     restapi.NewOptString(apiErr.Title),
				Detail:    restapi.NewOptString(apiErr.Detail),
				ErrorCode: restapi.NewOptString(fmt.Sprint(apiErr.ErrorCode)),
			},
		}
	case *ogenerrors.SecurityError:
		return &restapi.ProblemStatusCode{
			StatusCode: http.StatusForbidden,
			Response: restapi.Problem{
				Title: restapi.NewOptString(http.StatusText(http.StatusForbidden)),
			},
		}
	default:
		return &restapi.ProblemStatusCode{
			StatusCode: http.StatusInternalServerError,
			Response: restapi.Problem{
				Status: restapi.NewOptInt32(http.StatusInternalServerError),
				Title:  restapi.NewOptString(http.StatusText(http.StatusInternalServerError)),
			},
		}
	}
}
