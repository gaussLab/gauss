package simulation 

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/gauss/gauss/v4/x/orderbook/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyPoolMaxCount),
			func(r *rand.Rand) string {
				return fmt.Sprintf("%d", GenMaxPools(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyPoolMinPledgeAmount),
			func(r *rand.Rand) string {
				return fmt.Sprintf("%d", GenMinPledge(r))
			},
		),
	}
}
