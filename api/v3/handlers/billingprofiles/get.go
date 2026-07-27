package billingprofiles

import (
	"context"
	"errors"
	"net/http"

	api "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/api/v3/apierrors"
	"github.com/Pototoooo/meterforge/meterforge/billing"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
	"github.com/Pototoooo/meterforge/pkg/models"
)

type (
	GetBillingProfileRequest  = billing.GetProfileInput
	GetBillingProfileResponse = api.BillingProfile
	GetBillingProfileParams   = string
	GetBillingProfileHandler  httptransport.HandlerWithArgs[GetBillingProfileRequest, GetBillingProfileResponse, GetBillingProfileParams]
)

func (h *handler) GetBillingProfile() GetBillingProfileHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, profileID GetBillingProfileParams) (GetBillingProfileRequest, error) {
			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return GetBillingProfileRequest{}, err
			}

			return GetBillingProfileRequest{
				Profile: billing.ProfileID(models.NamespacedID{
					Namespace: ns,
					ID:        profileID,
				}),
				Expand: billing.ProfileExpand{
					Apps: true,
				},
			}, nil
		},
		func(ctx context.Context, request GetBillingProfileRequest) (GetBillingProfileResponse, error) {
			profile, err := h.service.GetProfile(ctx, request)
			if err != nil {
				return GetBillingProfileResponse{}, err
			}

			if profile == nil {
				return GetBillingProfileResponse{}, apierrors.NewNotFoundError(ctx, errors.New("billing profile not found"), "billing profile")
			}

			return ToAPIBillingProfile(*profile)
		},
		commonhttp.JSONResponseEncoderWithStatus[GetBillingProfileResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("get-billing-profile"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
			httptransport.WithErrorEncoder(errorEncoder()),
		)...,
	)
}
