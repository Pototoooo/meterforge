package httpdriver

import (
	"context"
	"net/http"

	"github.com/Pototoooo/meterforge/meterforge/notification"
	"github.com/Pototoooo/meterforge/meterforge/notification/webhook"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport/encoder"
	"github.com/Pototoooo/meterforge/pkg/models"
)

func errorEncoder() encoder.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter, r *http.Request) bool {
		return commonhttp.HandleErrorIfTypeMatches[notification.NotFoundError](ctx, http.StatusNotFound, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*feature.FeatureNotFoundError](ctx, http.StatusBadRequest, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*models.GenericValidationError](ctx, http.StatusBadRequest, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[webhook.ValidationError](ctx, http.StatusInternalServerError, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[webhook.NotFoundError](ctx, http.StatusInternalServerError, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[notification.UpdateAfterDeleteError](ctx, http.StatusConflict, err, w)
	}
}
