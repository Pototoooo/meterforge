package productcatalogdriver

import (
	"context"
	"net/http"

	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport/encoder"
	"github.com/Pototoooo/meterforge/pkg/pagination"
)

func getErrorEncoder() encoder.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter, r *http.Request) bool {
		return commonhttp.HandleErrorIfTypeMatches[*feature.FeatureNotFoundError](ctx, http.StatusNotFound, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*feature.FeatureInvalidFiltersError](ctx, http.StatusBadRequest, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*feature.ForbiddenError](ctx, http.StatusBadRequest, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*pagination.InvalidError](ctx, http.StatusBadRequest, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*feature.FeatureInvalidMeterAggregationError](ctx, http.StatusBadRequest, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*meter.MeterNotFoundError](ctx, http.StatusNotFound, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*feature.FeatureWithNameAlreadyExistsError](ctx, http.StatusConflict, err, w)
	}
}
