package chargeadapter

import (
	"github.com/samber/lo"

	chargecreditpurchase "github.com/Pototoooo/meterforge/meterforge/billing/charges/creditpurchase"
	chargeflatfee "github.com/Pototoooo/meterforge/meterforge/billing/charges/flatfee"
	"github.com/Pototoooo/meterforge/meterforge/billing/charges/meta"
	chargeusagebased "github.com/Pototoooo/meterforge/meterforge/billing/charges/usagebased"
	"github.com/Pototoooo/meterforge/meterforge/ledger"
	"github.com/Pototoooo/meterforge/pkg/models"
)

func chargeAnnotationsForCreditPurchaseCharge(charge chargecreditpurchase.Charge) models.Annotations {
	return chargeTransactionAnnotations(
		models.NamespacedID{
			Namespace: charge.Namespace,
			ID:        charge.ID,
		},
		charge.Intent.Subscription,
		nil,
	)
}

func chargeAnnotationsForFlatFeeCharge(charge chargeflatfee.Charge) models.Annotations {
	return chargeTransactionAnnotations(
		models.NamespacedID{
			Namespace: charge.Namespace,
			ID:        charge.ID,
		},
		charge.Intent.GetSubscription(),
		charge.State.FeatureID,
	)
}

func chargeAnnotationsForUsageBasedCharge(charge chargeusagebased.Charge) models.Annotations {
	return chargeTransactionAnnotations(
		models.NamespacedID{
			Namespace: charge.Namespace,
			ID:        charge.ID,
		},
		charge.Intent.GetSubscription(),
		lo.EmptyableToPtr(charge.State.FeatureID),
	)
}

func chargeTransactionAnnotations(chargeID models.NamespacedID, subscription *meta.SubscriptionReference, featureID *string) models.Annotations {
	var subscriptionID *string
	var subscriptionPhaseID *string
	var subscriptionItemID *string

	if subscription != nil {
		subscriptionID = &subscription.SubscriptionID
		subscriptionPhaseID = &subscription.PhaseID
		subscriptionItemID = &subscription.ItemID
	}

	return ledger.ChargeTransactionAnnotations(ledger.ChargeTransactionAnnotationsInput{
		ChargeID:            chargeID,
		SubscriptionID:      subscriptionID,
		SubscriptionPhaseID: subscriptionPhaseID,
		SubscriptionItemID:  subscriptionItemID,
		FeatureID:           featureID,
	})
}
