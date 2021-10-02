package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/gauss/gauss/v4/x/token/types"
)

// GetCmdQueryParams implements the query token related param command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query global token parameters.",
		Long:  strings.TrimSpace(
			fmt.Sprintf(`Query global token parameters

Example:
$ %s query %s params`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryTokens implements the query tokens command.
func GetCmdQueryTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens [owner]",
		Args:  cobra.RangeArgs(0,1),
		Short:  "Query token by the owner.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all tokens. optionall get a special token.


Example:
$ %s query %s tokens <owner>`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			var owner sdk.AccAddress
			if len(args) > 0 {
				owner, err = sdk.AccAddressFromBech32(args[0])
				if err != nil {
					return err
				}
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			res, err := queryClient.Tokens(
				context.Background(),
				&types.QueryTokensRequest{
					Owner:      owner.String(),
					Pagination: pageReq,
				},
			)
			if err != nil {
				return err
			}

			tokens := make([]types.TokenI, 0, len(res.Tokens))
			for _, eviAny := range res.Tokens {
				var evi types.TokenI
				if err = clientCtx.InterfaceRegistry.UnpackAny(eviAny, &evi); err != nil {
					return err
				}
				tokens = append(tokens, evi)
			}

			return clientCtx.PrintObjectLegacy(tokens)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "all tokens")

	return cmd
}

// GetCmdQueryToken implements the query token command.
func GetCmdQueryToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "token [symbol]",
		Args:  cobra.ExactArgs(1),
		Short: "Query a token by symbol.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a token by symbol

Example:
$ %s query %s token <symbol>`,
				version.AppName, types.ModuleName,
			),
	 	),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			if err := types.ValidateSymbol(args[0]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Token(
				context.Background(), 
				&types.QueryTokenRequest{
					Symbol: args[0],
				},
			)
			if err != nil {
				return err
			}

			// return clientCtx.PrintProto(res.Token)
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryTokenFees implements the query token fees command.
func GetCmdQueryTokenFees() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fees [symbol]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the token fees.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the token fees

Example:
$ %s query token fees <symbol>`,
				version.AppName, types.ModuleName,	
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			if err := types.ValidateSymbol(args[0]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Fees(
				context.Background(),
				&types.QueryFeesRequest{
					Symbol: args[0],
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd

}

// GetCmdQueryBurntoken implements the query burnt token command.
func GetCmdQueryBurntoken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burned [symbol]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the burned token.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the burned token

Example:
$ %s query %s burned <symbol>`,
				version.AppName, types.ModuleName,	
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			if err := types.ValidateSymbol(args[0]); err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Burntoken(
				context.Background(),
				&types.QueryBurntokenRequest{
					Symbol: args[0],
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}


