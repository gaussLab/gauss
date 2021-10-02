package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/gauss/gauss/v4/x/orderbook/types"
)

// NewTxCmd returns a root CLI command handler for all x/orderbook transaction commands.
func NewTxCmd() *cobra.Command {
	orderbookTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Orderbook transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	orderbookTxCmd.AddCommand(
		NewCreateTxPoolCmd(),
		NewAddPledgeCmd(),
		NewRedeemPledgeCmd(),
		NewPlaceOrderCmd(),
		NewRevokeOrderCmd(),
		NewAgreeOrderPairCmd(),
	)

	return orderbookTxCmd
}

func NewCreateTxPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [pledge]",
		Args: cobra.ExactArgs(1),
		Short: "create new tx pool.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`create new tx pool owning an address.

Example:
$ %s tx %s create-pool 100ugauss --from mykey 
`,
				version.AppName, types.ModuleName, 
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			ownerAddr := clientCtx.GetFromAddress()
			pledge, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}
			var defiAddr sdk.ValAddress = nil
			if defiAddress, _ := cmd.Flags().GetString(FlagDefiAddress); defiAddress != "" {
				defiAddr, err = sdk.ValAddressFromBech32(defiAddress)
				if err != nil {
					return err
				}
			}
			var delegatorAddr sdk.AccAddress = ownerAddr
			if delegatorAddress, _ := cmd.Flags().GetString(FlagDelegatorAddress); delegatorAddress != "" {
				delegatorAddr, err = sdk.AccAddressFromBech32(delegatorAddress)
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgCreatePool(ownerAddr, delegatorAddr, defiAddr, pledge)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsPools)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewAddPledgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-pledge [pledge]",
		Args: cobra.ExactArgs(1),
		Short: "add pledge to a pool.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`add pledge to a pool

Example:
$ %s tx %s add-pledge 100ugauss --from mykey
`,
				version.AppName, types.ModuleName, 
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			ownerAddr := clientCtx.GetFromAddress()
			pledge, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgAddPledge(ownerAddr, pledge)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewRedeemPledgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem-pledge [amount]",
		Args: cobra.ExactArgs(1),
		Short: "remove pool when redeeming all pledges.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`remove pool when redeeming all pledges

Example:
$ %s tx %s redeem-pledge 1000ugauss --from mykey
`,
				version.AppName, types.ModuleName, 
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			ownerAddr := clientCtx.GetFromAddress()
			amount, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgRedeemPledge(ownerAddr, amount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewPlaceOrderCmd() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "place-order [pool] [my-asset] [expect-asset] [price] [nonce]",
		Args: cobra.ExactArgs(5),
		Short: "place an order to a pool.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`place an order to a pool

Example:
$ %s tx %s place-order %s1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p 1000ugauss 5000uusdg 0.05uusdg 0001 --from mykey
`,
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			ownerAddr := clientCtx.GetFromAddress()

			poolAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			my_asset, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}
			expect_asset, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			// price, err := sdk.ParseCoinNormalized(args[3])
			price, err := sdk.ParseDecCoin(args[3])
			if err != nil {
				return err
			}
		
			order_id, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgPlaceOrder(poolAddr, ownerAddr, my_asset, expect_asset, price, order_id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewRevokeOrderCmd() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "revoke-order [pool] [tx-pair] [order_id]",
		Args: cobra.ExactArgs(3),
		Short: "revoke an order with order_id.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`revoke an order with order_id from special pool

Example:
$ %s tx %s revoke-order %s1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p gauss-usdg 10012 --is-left-order --from mykey
`,
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()

			poolAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			order_id, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			is_left_order, _ := cmd.Flags().GetBool(FlagIsLeftOrder)

			msg := types.NewMsgRevokeOrder(poolAddr, delAddr, args[1], is_left_order, order_id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsPools)
	cmd.Flags().AddFlagSet(fsOrders)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewAgreeOrderPairCmd() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "agree-orders [pool] [tx-pair] [left_order_id] [right_order_id] [price] [amount]",
		Args: cobra.RangeArgs(5,6),
		Short: "agree to transact with two orders in a pool.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`agree to transact with two orders in a pool

Example:
$ %s tx %s agree-order %s1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p gauss-usdg 10001 10002 0.01uusdg 20 --from mykey
`,
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			delAddr := clientCtx.GetFromAddress()

			poolAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			left_order_id, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			right_order_id, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}

			// price, err := sdk.ParseCoinNormalized(args[4])
			price, err := sdk.ParseDecCoin(args[4])
			if err != nil {
				return err
			}

			amount := sdk.NewInt(-1)
			if len(args) == 6 {
				amountL, ok := sdk.NewIntFromString(args[5])
				if !ok {
					return types.ErrInvalidAmount
				}
				if !amountL.IsPositive() {
					return types.ErrInvalidAmount
				}

				amount = amountL;
			}

			msg := types.NewMsgAgreeOrderPair(delAddr, poolAddr, args[1], left_order_id, right_order_id, price, amount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
