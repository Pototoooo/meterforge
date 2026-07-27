package apps

import (
	"context"
	"net/http"

	"github.com/Pototoooo/meterforge/meterforge/app"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport/encoder"
)

func errorEncoder() encoder.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter, r *http.Request) bool {
		return commonhttp.HandleErrorIfTypeMatches[*app.AppNotFoundError](ctx, http.StatusNotFound, err, w)
	}
}
