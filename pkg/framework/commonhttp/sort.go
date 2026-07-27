package commonhttp

import (
	"github.com/Pototoooo/meterforge/pkg/convert"
	"github.com/Pototoooo/meterforge/pkg/defaultx"
	"github.com/Pototoooo/meterforge/pkg/sortx"
)

func GetSortOrder[TInput comparable](asc TInput, inp *TInput) sortx.Order {
	return defaultx.WithDefault(
		convert.SafeDeRef[TInput, sortx.Order](
			inp,
			func(o TInput) *sortx.Order {
				if o == asc {
					return convert.ToPointer(sortx.OrderAsc)
				}
				return convert.ToPointer(sortx.OrderDesc)
			},
		),
		sortx.OrderNone,
	)
}
