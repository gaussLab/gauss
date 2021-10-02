package cli

import (
	flag "github.com/spf13/pflag"

	"github.com/gauss/gauss/v4/x/token/types"
)

const (
	FlagName 	  = "name"
	FlagSymbol        = "symbol"
	FlagDecimals      = "decimals"
	FlagSmallestUnit  = "smallest-unit"
	FlagInitialSupply = "initial-supply"
	FlagTotalSupply   = "total-supply"
	FlagMintable      = "mintable"
	FlagUnlocked      = "unlocked"
	FlagTo            = "to"
	FlagAmount        = "amount"
)

var (
	FsIssueToken         = flag.NewFlagSet("", flag.ContinueOnError)
	FsEditToken          = flag.NewFlagSet("", flag.ContinueOnError)
	FsMintToken          = flag.NewFlagSet("", flag.ContinueOnError)
	FsTransferTokenOwner = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsIssueToken.String(FlagName, "", "The token name. Once created, it cannot be modified")
	// FsIssueToken.String(FlagSymbol, "", "The token symbol. Once created, it cannot be modified")
	FsIssueToken.Uint32(FlagDecimals, 6, "The token decimals. The maximum value is 18")
	FsIssueToken.String(FlagSmallestUnit, "", "The token the smallest unit")
	FsIssueToken.Uint64(FlagInitialSupply, 0, "The initial supply of the token. The unit is smallest-unit")
	FsIssueToken.Uint64(FlagTotalSupply, types.MaximumAmount, "The maximum supply of the token(default 1.8*10^19). The unit is smallest-unit")
	FsIssueToken.Bool(FlagMintable, false, "Whether the token can be minted (default false)")
	FsIssueToken.Bool(FlagUnlocked, true, "Whether the token can be transfer")

	FsEditToken.Bool(FlagMintable, false, "Whether the token can be minted, default to false")

	FsMintToken.String(FlagTo, "", "Address to which the token is to be minted")
	FsMintToken.Uint64(FlagAmount, 0, "Amount of the token to be minted")

	FsTransferTokenOwner.String(FlagTo, "", "The new owner")
}
