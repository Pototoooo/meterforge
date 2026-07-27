package billingprofiles

import (
	"context"
	"errors"
	"net/http"

	api "github.com/Pototoooo/meterforge/api/v3"
	"github.com/Pototoooo/meterforge/api/v3/apierrors"
	"github.com/Pototoooo/meterforge/api/v3/request"
	"github.com/Pototoooo/meterforge/meterforge/billing"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport"
	"github.com/Pototoooo/meterforge/pkg/models"
)

type (
	UpdateBillingProfileRequest struct {
		Namespace string
		ProfileID string
		Profile   billing.UpdateProfileInput
	}
	UpdateBillingProfileParams   = string
	UpdateBillingProfileResponse = api.BillingProfile
	UpdateBillingProfileHandler  httptransport.HandlerWithArgs[UpdateBillingProfileRequest, UpdateBillingProfileResponse, UpdateBillingProfileParams]
)

func (h *handler) UpdateBillingProfile() UpdateBillingProfileHandler {
	return httptransport.NewHandlerWithArgs(
		func(ctx context.Context, r *http.Request, profileID UpdateBillingProfileParams) (UpdateBillingProfileRequest, error) {
			body := api.UpsertBillingProfileRequest{}
			if err := request.ParseBody(r, &body); err != nil {
				return UpdateBillingProfileRequest{}, err
			}

			ns, err := h.resolveNamespace(ctx)
			if err != nil {
				return UpdateBillingProfileRequest{}, err
			}

			profile, err := FromAPIUpsertBillingProfileRequest(models.NamespacedID{
				Namespace: ns,
				ID:        profileID,
			}, body)
			if err != nil {
				return UpdateBillingProfileRequest{}, err
			}

			req := UpdateBillingProfileRequest{
				Namespace: ns,
				ProfileID: profileID,
				Profile:   profile,
			}

			return req, nil
		},
		func(ctx context.Context, request UpdateBillingProfileRequest) (UpdateBillingProfileResponse, error) {
			updatedProfile, err := h.service.UpdateProfile(ctx, request.Profile)
			if err != nil {
				return UpdateBillingProfileResponse{}, err
			}

			if updatedProfile == nil {
				return UpdateBillingProfileResponse{}, errors.New("failed to update billing profile")
			}

			return ToAPIBillingProfile(*updatedProfile)
		},
		commonhttp.JSONResponseEncoderWithStatus[UpdateBillingProfileResponse](http.StatusOK),
		httptransport.AppendOptions(
			h.options,
			httptransport.WithOperationName("update-billing-profile"),
			httptransport.WithErrorEncoder(apierrors.GenericErrorEncoder()),
			httptransport.WithErrorEncoder(errorEncoder()),
		)...,
	)
}
