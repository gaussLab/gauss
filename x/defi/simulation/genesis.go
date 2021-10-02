package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/gauss/gauss/v4/x/defi/types"
)

// Simulation parameter constants
const (
	mintInflationKey  = "mint_inflation"
	communityTaxKey   = "community_tax"
	commissionRateKey = "commission_rate"
	marketRateKey     = "market_rate"
	unbondingTimeKey  = "unbonding_time"
	maxDefisKey       = "max_defis"
	maxEntriesKey     = "max_entries"
	historicalEntriesKey = "historical_entries"
)

// GenMintInflation randomized MintInflation
func GenMintInflation(r *rand.Rand) (coin sdk.Coin) {
	return sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(1000)}
}

// GenCommunityTax randomized CommunityTax
func GenCommunityTax(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(1, 2).Add(sdk.NewDecWithPrec(int64(r.Intn(30)), 2))
}

// GenCommissionRate randomized CommissionRate
func GenCommissionRate(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(1, 2).Add(sdk.NewDecWithPrec(int64(r.Intn(30)), 2))
}

// GenMarketRate randomized MarketRate
func GenMarketRate(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(1, 2).Add(sdk.NewDecWithPrec(int64(r.Intn(30)), 2))
}

// GenUnbondingTime randomized UnbondingTime
func GenUnbondingTime(r *rand.Rand) (ubdTime time.Duration) {
	return time.Duration(simulation.RandIntBetween(r, 60, 60*60*24*3*2)) * time.Second
}

// GenMaxDefis randomized MaxDefis
func GenMaxDefis(r *rand.Rand) (maxDefis uint32) {
	return uint32(r.Intn(250) + 1)
}

// GenMaxEntries randomized MaxEntries
func GenMaxEntries(r *rand.Rand) (maxEntries uint32) {
	return uint32(r.Intn(250) + 1)
}

// GetHistEntries randomized HistoricalEntries between 0-100.
func GetHistEntries(r *rand.Rand) uint32 {
	return uint32(r.Intn(int(types.DefaultHistoricalEntries + 1)))
}

// RandomizedGenState generates a random GenesisState for staking
func RandomizedGenState(simState *module.SimulationState) {
	// params
	var (
		mintInflation sdk.Coin
		communityTax  sdk.Dec
		commissionRate sdk.Dec
		marketRate    sdk.Dec
		unbondTime    time.Duration
		maxDefis      uint32
		maxEntries    uint32
		histEntries   uint32
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, mintInflationKey, &mintInflation, simState.Rand,
		func(r *rand.Rand) { mintInflation = GenMintInflation(r) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, communityTaxKey, &communityTax, simState.Rand,
		func(r *rand.Rand) { communityTax = GenCommunityTax(r) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, commissionRateKey, &commissionRate, simState.Rand,
		func(r *rand.Rand) { commissionRate = GenCommissionRate(r) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, marketRateKey, &marketRate, simState.Rand,
		func(r *rand.Rand) { marketRate = GenMarketRate(r) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, maxDefisKey, &maxDefis, simState.Rand,
		func(r *rand.Rand) { maxDefis = GenMaxDefis(r) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, maxEntriesKey, &maxEntries, simState.Rand,
		func(r *rand.Rand) { maxEntries = GenMaxEntries(r) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, historicalEntriesKey, &histEntries, simState.Rand,
		func(r *rand.Rand) { histEntries = GetHistEntries(r) },
	)

	
	// NOTE: the slashing module need to be defined after the staking module on the
	// NewSimulationManager constructor for this to work
	simState.UnbondTime = unbondTime
	params := types.NewParams(sdk.DefaultBondDenom, mintInflation, communityTax, commissionRate, marketRate,
		 simState.UnbondTime, maxDefis, maxEntries, histEntries)

	// defis & delegations
	var (
		defis  []types.Defi
		delegations []types.Delegation
	)

	defiAddrs := make([]sdk.ValAddress, simState.NumBonded)

	for i := 0; i < int(simState.NumBonded); i++ {
		defiAddr := sdk.ValAddress(simState.Accounts[i].Address)
		defiAddrs[i] = defiAddr

		defi, err := types.NewDefi(defiAddr, types.Description{})
		if err != nil {
			panic(err)
		}
		defi.Tokens = sdk.NewInt(simState.InitialStake)
		defi.DelegatorShares = sdk.NewDec(simState.InitialStake)

		delegation := types.NewDelegation(
			simState.Accounts[i].Address,
			defiAddr,
			sdk.NewDec(simState.InitialStake),
		)

		defis = append(defis, defi)
		delegations = append(delegations, delegation)
	}

	stakingGenesis := types.GenesisState{
		FeePool: types.InitialFeePool(),
		Params: params,
		Defis: defis,
		Delegations: delegations,
	}

	bz, err := json.MarshalIndent(&stakingGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated staking parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&stakingGenesis)
}
