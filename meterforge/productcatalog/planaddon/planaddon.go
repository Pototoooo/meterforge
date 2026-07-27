package planaddon

import (
	"errors"

	"github.com/Pototoooo/meterforge/meterforge/productcatalog"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/addon"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog/plan"
	"github.com/Pototoooo/meterforge/pkg/models"
)

var _ models.Validator = (*PlanAddon)(nil)

type PlanAddon struct {
	models.NamespacedID
	models.ManagedModel

	productcatalog.PlanAddonMeta

	// Addon
	Plan plan.Plan `json:"plan"`

	// Addon
	Addon addon.Addon `json:"addon"`
}

func (a PlanAddon) Validate() error {
	var errs []error

	if err := a.NamespacedID.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := a.ManagedModel.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := a.Plan.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := a.Addon.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := a.AsProductCatalogPlanAddon().Validate(); err != nil {
		errs = append(errs, err)
	}

	return models.NewNillableGenericValidationError(errors.Join(errs...))
}

func (a PlanAddon) AsProductCatalogPlanAddon() productcatalog.PlanAddon {
	return productcatalog.PlanAddon{
		PlanAddonMeta: a.PlanAddonMeta,
		Plan:          a.Plan.AsProductCatalogPlan(),
		Addon:         a.Addon.AsProductCatalogAddon(),
	}
}
