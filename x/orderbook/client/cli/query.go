package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/gauss/gauss/v4/x/orderbook/types"

)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	orderbookQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the orderbook module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	orderbookQueryCmd.AddCommand(
		GetCmdQueryTxPairStats(),
		GetCmdQueryOrders(),
		GetCmdQueryTxPairs(),
		GetCmdQueryPools(),
		GetCmdQueryParams(),
	)

	return orderbookQueryCmd
}

// GetCmdQueryTxPairs implements the tx-pairs query command.
func GetCmdQueryTxPairStats() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "tx-pair-stats [pool]",
		Args: cobra.ExactArgs(1),
		Short: "Query an tx-pair from a pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a tx-pair from a pool.

Example:
$ %s query %s tx-pair-stats %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
`,
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			poolAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryTxPairsStatsRequest{PoolAddress: poolAddr.String(), Pagination: pageReq}
			res, err := queryClient.TxPairsStats(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "txPairStats")

	return cmd
}

// GetCmdQueryOrders implements the orders query command.
func GetCmdQueryOrders() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "orders [pool] [tx-pair] [order_id]",
		Args: cobra.RangeArgs(2,3),
		Short: "Query all orders",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all orders. optionally restrict to a single order

Example:
$ %s query %s orders %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj ugauss/uusdg 
$ %s query %s orders %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj ugauss/uusdg 1001 --is-left-order
`,
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			poolAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			isLeftOrder, _ := cmd.Flags().GetBool(FlagIsLeftOrder)

			if len(args) == 3 {
				orderId, err := strconv.ParseUint(args[2], 10, 64)
				if err != nil {
					return err
				}

				msg := &types.QueryPoolTxPairOrderRequest{
					PoolAddress: poolAddr.String(), 
					TxPair: args[1], 
					OrderId: orderId, 
					IsLeftOrder: isLeftOrder}
				res, err := queryClient.Order(context.Background(), msg)
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res)
			}

			msg := &types.QueryPoolTxPairOrdersRequest{PoolAddress: poolAddr.String(), TxPair: args[1], IsLeftOrder: isLeftOrder}
			res, err := queryClient.Orders(context.Background(), msg)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().AddFlagSet(fsOrders)
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "orders")

	return cmd
}

// GetCmdQueryTxPairs implements the txPairs query command.
func GetCmdQueryTxPairs() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "tx-pairs [pool] [tx-pair]",
		Args: cobra.RangeArgs(1, 2),
		Short: "Query all tx-Pairs",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all tx-Pairs. optionally restrict to a single txPair

Example:
$ %s query %s tx-pairs %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj
$ %s query %s tx-pairs %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj ugauss/uusdg
`,
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			poolAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			if len(args) == 2 {
				res, err := queryClient.TxPair(context.Background(), 
					&types.QueryPoolTxPairRequest{PoolAddress: poolAddr.String(), TxPair: args[1]})
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res)
			}

			res, err := queryClient.TxPairs(context.Background(), &types.QueryPoolTxPairsRequest{PoolAddress: poolAddr.String()})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "txPairs")

	return cmd
}

// GetCmdQueryPools implements the pools query command.
func GetCmdQueryPools() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "pools [pool]",
		Args: cobra.RangeArgs(0, 1),
		Short: "Query all pools",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all pools. optionally restrict to a single pool

Example:
$ %s query %s pools 
$ %s query %s pools %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj 
`,
				version.AppName, types.ModuleName, 
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			if len(args) == 1 {
				poolAddr, err := sdk.AccAddressFromBech32(args[0])
				if err != nil {
					return err
				}
				
				res, err := queryClient.GetPool(context.Background(), 
					&types.QueryPoolRequest{PoolAddress: poolAddr.String()})
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res)
			}

			res, err := queryClient.Pools(context.Background(), &types.QueryPoolsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryParams implements the params query command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args: cobra.NoArgs,
		Short: "Query the orderbook module parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the orderbook module parameters information

Example:
$ %s query %s params 
`,
				version.AppName, types.ModuleName, 
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(),
				&types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
