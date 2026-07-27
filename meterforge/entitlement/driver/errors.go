package entitlementdriver

import (
	"context"
	"net/http"

	"github.com/Pototoooo/meterforge/meterforge/entitlement"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport/encoder"
	"github.com/Pototoooo/meterforge/pkg/pagination"
)

func GetErrorEncoder() encoder.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter, r *http.Request) bool {
		// user errors
		return commonhttp.HandleErrorIfTypeMatches[*feature.FeatureNotFoundError](ctx, http.StatusNotFound, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*pagination.InvalidError](ctx, http.StatusBadRequest, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*entitlement.NotFoundError](ctx, http.StatusNotFound, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*entitlement.AlreadyExistsError](
				ctx, http.StatusConflict, err, w,
				func(specificErr *entitlement.AlreadyExistsError) map[string]interface{} {
					return map[string]interface{}{
						"conflictingEntityId": specificErr.EntitlementID,
					}
				}) ||
			commonhttp.HandleErrorIfTypeMatches[*entitlement.InvalidValueError](ctx, http.StatusBadRequest, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*entitlement.InvalidFeatureError](ctx, http.StatusBadRequest, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*entitlement.WrongTypeError](ctx, http.StatusBadRequest, err, w)
	}
}
