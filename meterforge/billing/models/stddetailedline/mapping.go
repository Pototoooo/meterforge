package stddetailedline

import (
	"time"

	"github.com/alpacahq/alpacadecimal"
	"github.com/samber/lo"

	"github.com/Pototoooo/meterforge/meterforge/billing/models/creditsapplied"
	"github.com/Pototoooo/meterforge/meterforge/billing/models/externalid"
	"github.com/Pototoooo/meterforge/meterforge/billing/models/totals"
	"github.com/Pototoooo/meterforge/meterforge/productcatalog"
	"github.com/Pototoooo/meterforge/pkg/convert"
	"github.com/Pototoooo/meterforge/pkg/models"
	"github.com/Pototoooo/meterforge/pkg/timeutil"
)

type DBGetter interface {
	GetNamespace() string
	GetID() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetDeletedAt() *time.Time
	GetName() string
	GetDescription() *string
	GetCategory() Category
	GetChildUniqueReferenceID() string
	GetIndex() *int
	GetPaymentTerm() productcatalog.PaymentTermType
	GetServicePeriodStart() time.Time
	GetServicePeriodEnd() time.Time
	GetPerUnitAmount() alpacadecimal.Decimal
	GetQuantity() alpacadecimal.Decimal
	GetCreditsApplied() *creditsapplied.CreditsApplied

	externalid.LineExternalIDGetter
	totals.TotalsGetter
}

func FromDB[T DBGetter](dbEntity T) Base {
	return Base{
		ManagedResource: models.NewManagedResource(models.ManagedResourceInput{
			Namespace:   dbEntity.GetNamespace(),
			ID:          dbEntity.GetID(),
			CreatedAt:   dbEntity.GetCreatedAt().In(time.UTC),
			UpdatedAt:   dbEntity.GetUpdatedAt().In(time.UTC),
			DeletedAt:   convert.TimePtrIn(dbEntity.GetDeletedAt(), time.UTC),
			Name:        dbEntity.GetName(),
			Description: dbEntity.GetDescription(),
		}),
		Category:               dbEntity.GetCategory(),
		ChildUniqueReferenceID: dbEntity.GetChildUniqueReferenceID(),
		Index:                  dbEntity.GetIndex(),
		PaymentTerm:            dbEntity.GetPaymentTerm(),
		ServicePeriod: timeutil.ClosedPeriod{
			From: dbEntity.GetServicePeriodStart().In(time.UTC),
			To:   dbEntity.GetServicePeriodEnd().In(time.UTC),
		},
		PerUnitAmount:  dbEntity.GetPerUnitAmount(),
		Quantity:       dbEntity.GetQuantity(),
		Totals:         totals.FromDB(dbEntity),
		ExternalIDs:    externalid.MapLineExternalIDFromDB(dbEntity),
		CreditsApplied: lo.FromPtr(dbEntity.GetCreditsApplied()),
	}
}
