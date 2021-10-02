package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/gauss/gauss/v4/x/token/types"
)

// GetCmdIssueToken implements the issue token command
func GetCmdIssueToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "issue [symbol]",
		Args:  cobra.ExactArgs(1),
		Short: "Issue a new token.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Issue a new token 

Example:
$ %s tx %s issue gauss --name=\"gauss network\" --smallest-unit=ugauss --decimals=6 --initial-supply=1000000000 --total-supply=10000000000 --mintable=true --from=<key-name>
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			symbol := args[0]
			name, err := cmd.Flags().GetString(FlagName)
			if err != nil {
				return err
			}
			// symbol, err := cmd.Flags().GetString(FlagSymbol)
			// if err != nil {
			//	return err
			// }
			smallestUnit, err := cmd.Flags().GetString(FlagSmallestUnit)
			if err != nil {
				return err
			}
			decimals, err := cmd.Flags().GetUint32(FlagDecimals)
			if err != nil {
				return err
			}
			initialSupply, err := cmd.Flags().GetUint64(FlagInitialSupply)
			if err != nil {
				return err
			}
			totalSupply, err := cmd.Flags().GetUint64(FlagTotalSupply)
			if err != nil {
				return err
			}
			mintable, err := cmd.Flags().GetBool(FlagMintable)
			if err != nil {
				return err
			}
			unlocked, err := cmd.Flags().GetBool(FlagUnlocked)
			if err != nil {
				return err
			}
			owner := clientCtx.GetFromAddress()

			msg := types.NewMsgIssueToken(name, symbol, smallestUnit, decimals, initialSupply, 
				totalSupply, mintable, unlocked, owner.String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			var prompt = "The token issuance transaction will consume extra fee"
			generateOnly, err := cmd.Flags().GetBool(flags.FlagGenerateOnly)
			if err != nil {
				return err
			}
			if !generateOnly {
				fees, errL := queryTokenFees(clientCtx, msg.Symbol)
				if errL != nil {
					return fmt.Errorf("failed to query token issuance fees: %s", errL.Error())
				}
				issueFees := sdk.Coins{fees.IssueFee}.String()
				prompt += fmt.Sprintf(": %s", issueFees)
			}
			
			fmt.Println(prompt)
			
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsIssueToken)
	_ = cmd.MarkFlagRequired(FlagName)
	// _ = cmd.MarkFlagRequired(FlagSymbol)
	_ = cmd.MarkFlagRequired(FlagSmallestUnit)
	_ = cmd.MarkFlagRequired(FlagInitialSupply)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdEditToken implements the edit token command
func GetCmdEditToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "edit [symbol]",
		Args: cobra.ExactArgs(1),
		Short: "Edit an existing token.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Edit an existing token

Example:
$ %s tx %s edit gauss --mintable=true --from=<key-name>
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			mintable, err := cmd.Flags().GetBool(FlagMintable)
			if err != nil {
				return err
			}
			owner := clientCtx.GetFromAddress()

			msg := types.NewMsgEditToken(args[0], mintable, owner.String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsEditToken)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdMintToken() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:  "mint [symbol]",
		Args: cobra.ExactArgs(1),
		Short: "Mint tokens to a specified address.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Mint tokens to a specified address

Example:
$ %s tx %s mint gauss --amount=1000ugauss --to=%s1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p --from=my_key
`,
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			owner := clientCtx.GetFromAddress()

			amount, err := cmd.Flags().GetUint64(FlagAmount)
			if err != nil {
				return err
			}

			toAddr, err := cmd.Flags().GetString(FlagTo)
			if err != nil {
				return err
			}
			if toAddr == "" {
				toAddr = owner.String()	
			}

			if _, err = sdk.AccAddressFromBech32(toAddr); err != nil {
				return err
			}

			msg := types.NewMsgMintToken(
				args[0], owner.String(), toAddr, amount,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			var prompt = "The token minting transaction will consume extra fee"
			generateOnly, err := cmd.Flags().GetBool(flags.FlagGenerateOnly) 
			if err != nil {
				return err
			}
			if !generateOnly {
				fees, errL := queryTokenFees(clientCtx, args[0])
				if errL != nil {
					return fmt.Errorf("failed to query token minting fees: %s", errL.Error())
				}
				mintFees := sdk.Coins{fees.MintFee}.String()
				prompt += fmt.Sprintf(": %s", mintFees)
			}

			fmt.Println(prompt)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsMintToken)
	_ = cmd.MarkFlagRequired(FlagAmount)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetCmdBurnToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "burn [symbol]",
		Args: cobra.ExactArgs(1),
		Short: "Burn amount of token.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Burn amount of token

Example:
$ %s tx %s burn gauss --amount=1000ugauss --from=my_key
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := cmd.Flags().GetUint64(FlagAmount)
			if err != nil {
				return err
			}
			owner := clientCtx.GetFromAddress()

			msg := types.NewMsgBurnToken(
				args[0], owner.String(), amount,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsMintToken)
	_ = cmd.MarkFlagRequired(FlagAmount)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdUnlockToken implements unlock the token  command
func GetCmdUnlockToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "unlock [symbol]",
		Args: cobra.ExactArgs(1),
		Short: "Unlock the locked token.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Unlock the locked token

Example:
$ %s tx %s unlock gauss --from=my_key
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			owner := clientCtx.GetFromAddress()

			msg := types.NewMsgUnlockToken(args[0], owner.String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdTransferTokenOwner implements transfer the token owner command
func GetCmdTransferTokenOwner() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:  "transfer [symbol]",
		Args: cobra.ExactArgs(1),
		Short: "Transfer the owner of a token to a new owner.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Transfer the owner of a token to a new owner

Example:
$ %s tx %s transfer gauss --to=%s1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p  --from=my_key
`,
				version.AppName, types.ModuleName, bech32PrefixAccAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			toAddr, err := cmd.Flags().GetString(FlagTo)
			if err != nil {
				return err
			}
			if _, err := sdk.AccAddressFromBech32(toAddr); err != nil {
				return err
			}
			owner := clientCtx.GetFromAddress()

			msg := types.NewMsgTransferTokenOwner(args[0], owner.String(), toAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsTransferTokenOwner)
	_ = cmd.MarkFlagRequired(FlagTo)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
