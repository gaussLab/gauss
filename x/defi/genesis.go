package defi

import (
	"fmt"

	"github.com/gauss/gauss/v4/x/defi/types"
)

// ValidateGenesis validates the provided staking genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate defis)
func ValidateGenesis(data *types.GenesisState) error {
	if err := validateGenesisStateDefis(data.Defis); err != nil {
		return err
	}
	if err := data.FeePool.ValidateGenesis(); err != nil {
		return err
	}

	return data.Params.Validate()
}

func validateGenesisStateDefis(defis []types.Defi) error {
	for i := 0; i < len(defis); i++ {
		defi := defis[i]
		if defi.DelegatorShares.IsZero() && !defi.IsUnbonding() {
			return fmt.Errorf("bonded/unbonded genesis defi cannot have zero delegator shares, defi: %v", defi)
		}
	}

	return nil
}
