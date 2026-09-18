package httpdriver

import (
	"context"
	"net/http"

	"github.com/Pototoooo/meterforge/meterforge/entitlement"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/feature"
	"github.com/Pototoooo/meterforge/meterforge/subscription"
	subscriptionentitlement "github.com/Pototoooo/meterforge/meterforge/subscription/entitlement"
	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/framework/transport/httptransport/encoder"
	"github.com/Pototoooo/meterforge/pkg/models"
	"github.com/Pototoooo/meterforge/pkg/pagination"
)

func errorEncoder() encoder.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter, r *http.Request) bool {
		issues, issueErr := models.AsValidationIssues(err)
		if issueErr == nil && len(issues) > 0 {
			mappedIssues, err := mapValidationIssueForAPI(issues)
			if err != nil {
				return false // Server dies if mapping fails
			}

			if commonhttp.HandleIssueIfHTTPStatusKnown(ctx, mappedIssues.AsError(), w) {
				return true
			}
		}

		// Generic errors
		return commonhttp.HandleErrorIfTypeMatches[*pagination.InvalidError](ctx, http.StatusBadRequest, err, w) ||
			// Subscription errors
			commonhttp.HandleErrorIfTypeMatches[*subscription.PatchConflictError](ctx, http.StatusConflict, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*subscription.PatchForbiddenError](ctx, http.StatusForbidden, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*subscription.PatchValidationError](ctx, http.StatusBadRequest, err, w) ||
			commonhttp.HandleErrorIfTypeMatches[*subscriptionentitlement.NotFoundError](ctx, http.StatusNotFound, err, w) ||
			// dependency: entitlement
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
			commonhttp.HandleErrorIfTypeMatches[*entitlement.WrongTypeError](ctx, http.StatusBadRequest, err, w) ||
			// dependency: plan
			commonhttp.HandleErrorIfTypeMatches[*feature.FeatureNotFoundError](ctx, http.StatusBadRequest, err, w)
	}
}

func mapValidationIssueForAPI(issues models.ValidationIssues) (models.ValidationIssues, error) {
	res := make(models.ValidationIssues, 0, len(issues))

	for _, issue := range issues {
		mapped, err := subscription.MapSubscriptionSpecValidationIssueField(issue)
		if err != nil {
			return res, err
		}

		res = append(res, mapped)
	}

	return res, nil
}
