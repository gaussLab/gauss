package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/gauss/gauss/v4/x/token/types"
)

// NewTxCmd returns the transaction commands for the token module.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Token transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		GetCmdIssueToken(),
		GetCmdEditToken(),
		GetCmdMintToken(),
		GetCmdBurnToken(),
		GetCmdUnlockToken(),
		GetCmdTransferTokenOwner(),
	)

	return txCmd
}


// GetQueryCmd returns the query commands for the token module.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the token module",
		DisableFlagParsing:         true,
                SuggestionsMinimumDistance: 2,
                RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryTokens(),
		GetCmdQueryToken(),
		GetCmdQueryTokenFees(),
		GetCmdQueryBurntoken(),
	)

	return queryCmd
}
