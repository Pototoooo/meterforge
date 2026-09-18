package payment

import (
	"net/http"

	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/models"
)

const ErrCodePaymentAlreadyAuthorized models.ErrorCode = "payment_already_authorized"

var ErrPaymentAlreadyAuthorized = models.NewValidationIssue(
	ErrCodePaymentAlreadyAuthorized,
	"payment already authorized",
	models.WithCriticalSeverity(),
	commonhttp.WithHTTPStatusCodeAttribute(http.StatusBadRequest),
)

const ErrCodePaymentAlreadySettled models.ErrorCode = "payment_already_settled"

var ErrPaymentAlreadySettled = models.NewValidationIssue(
	ErrCodePaymentAlreadySettled,
	"payment already settled",
	models.WithCriticalSeverity(),
	commonhttp.WithHTTPStatusCodeAttribute(http.StatusBadRequest),
)

const ErrCodeCannotSettleNotAuthorizedPayment models.ErrorCode = "cannot_settle_not_authorized_payment"

var ErrCannotSettleNotAuthorizedPayment = models.NewValidationIssue(
	ErrCodeCannotSettleNotAuthorizedPayment,
	"cannot settle an unauthorized payment",
	models.WithCriticalSeverity(),
	commonhttp.WithHTTPStatusCodeAttribute(http.StatusBadRequest),
)
