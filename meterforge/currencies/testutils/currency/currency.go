package currency

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Pototoooo/meterforge/meterforge/currencies"
	"github.com/Pototoooo/meterforge/pkg/currencyx"
)

func NewFiatCurrency[T ~string](t testing.TB, code T) currencies.Currency {
	t.Helper()

	currency, err := currencies.NewFiatCurrency(currencyx.Code(code))
	require.NoError(t, err)

	return currency
}
