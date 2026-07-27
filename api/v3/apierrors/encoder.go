package apierrors

import (
	"context"
	"net/http"

	"github.com/samber/lo"

	"github.com/Pototoooo/meterforge/meterforge/meter"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport/encoder"
)

// GenericErrorEncoder is an error encoder that encodes the error as a generic error.
func GenericErrorEncoder() encoder.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter, r *http.Request) bool {
		if err, ok := lo.ErrorsAs[*BaseAPIError](err); ok {
			err.HandleAPIError(w, r)
			return true
		}

		return commonhttp.HandleErrorIfTypeMatches[*feature.FeatureNotFoundError](ctx, http.StatusNotFound, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*meter.MeterNotFoundError](ctx, http.StatusNotFound, err, w) ||
			commonhttp.HandleIssueIfHTTPStatusKnown(ctx, err, w)
	}
}
