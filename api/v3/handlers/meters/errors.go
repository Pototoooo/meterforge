package meters

import (
	"net/http"

	"github.com/Pototoooo/meterforge/pkg/framework/commonhttp"
	"github.com/Pototoooo/meterforge/pkg/models"
)

const ErrCodeReservedDimension models.ErrorCode = "reserved_dimension"

var ErrReservedDimension = models.NewValidationIssue(
	ErrCodeReservedDimension,
	"dimension name is reserved",
	models.WithFieldString("dimensions"),
	models.WithCriticalSeverity(),
	commonhttp.WithHTTPStatusCodeAttribute(http.StatusBadRequest),
)

func NewReservedDimensionError(dimension string) error {
	return ErrReservedDimension.WithPathString("dimensions", dimension).WithAttr("value", dimension)
}
