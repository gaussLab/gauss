package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gauss/gauss/v4/x/orderbook/types"
)

// Simulation parameter constants
const (
	maxPools     = "max_pools"
	minPledge    = "min_pledge"
)

//  randomized maxPools
func GenMaxPools(r *rand.Rand) (maxPools uint32) {
	return uint32(r.Intn(250) + 1)
}

// randomized minPledge
func GenMinPledge(r *rand.Rand) uint32 {
	return uint32(r.Intn(types.DefaultPoolMinPledgeAmount.Sign()) + 1)
}

// RandomizedGenState generates a random GenesisState for staking
func RandomizedGenState(simState *module.SimulationState) {
	// params
	var (
		maxPoolsL     uint32
		minPledgeL    uint32
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, maxPools, &maxPoolsL, simState.Rand,
		func(r *rand.Rand) { maxPoolsL = GenMaxPools(r) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, minPledge, &minPledgeL, simState.Rand,
		func(r *rand.Rand) { minPledgeL = GenMinPledge(r) },
	)

	params := types.DefaultParams()
	orderbookGenesis := types.NewGenesisState(params, []types.Pool{}, []types.Order{}, []types.TxPairStats{})

	bz, err := json.MarshalIndent(&orderbookGenesis.Params, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated orderbook parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(orderbookGenesis)
}
