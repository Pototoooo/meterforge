package billing

import (
	"github.com/spf13/cobra"

	"github.com/Pototoooo/meterforge/cmd/jobs/billing/advance"
	"github.com/Pototoooo/meterforge/cmd/jobs/billing/advancecharges"
	"github.com/Pototoooo/meterforge/cmd/jobs/billing/collect"
	"github.com/Pototoooo/meterforge/cmd/jobs/billing/subscriptionsync"
)

var Cmd = &cobra.Command{
	Use:   "billing",
	Short: "Billing operations",
}

func init() {
	Cmd.AddCommand(advance.Cmd)
	Cmd.AddCommand(advancecharges.Cmd)
	Cmd.AddCommand(collect.Cmd)
	Cmd.AddCommand(subscriptionsync.Cmd)
}
