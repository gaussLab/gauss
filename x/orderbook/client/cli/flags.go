package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagDefiAddress    = "defi-address"
	FlagDelegatorAddress    = "delegator-address"
	FlagIsLeftOrder    = "is-left-order"
)

var (
	fsPools		= flag.NewFlagSet("", flag.ContinueOnError)
	fsOrders	= flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsPools.String(FlagDefiAddress, "", "belong to defi enviroment")
	fsPools.String(FlagDelegatorAddress, "", "allow delegator to execute orders")
	fsOrders.Bool(FlagIsLeftOrder, false, "left order is buy order, right order is sale order. (default sale order)")
}
