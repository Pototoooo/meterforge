package features

import (
	"context"
	"net/http"

	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport/encoder"
)

func errorEncoder() encoder.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter, r *http.Request) bool {
		return commonhttp.HandleErrorIfTypeMatches[*feature.FeatureInvalidFiltersError](ctx, http.StatusBadRequest, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*feature.FeatureInvalidMeterAggregationError](ctx, http.StatusBadRequest, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*feature.ForbiddenError](ctx, http.StatusForbidden, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*feature.FeatureWithNameAlreadyExistsError](ctx, http.StatusConflict, err, w)
	}
}
