package orderbook

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gauss/gauss/v4/x/orderbook/keeper"
	"github.com/gauss/gauss/v4/x/orderbook/types"
)

// InitGenesis sets the pool and parameters for the provided keeper. 
func InitGenesis(
	ctx sdk.Context, keeper keeper.Keeper, accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper, data *types.GenesisState,
) (res []abci.ValidatorUpdate) {
	if err := ValidateGenesis(data); err != nil {
                panic(err.Error())
        }

	keeper.SetParams(ctx, data.Params)

	for _, pool := range data.Pools {
		poolAddr, err := sdk.AccAddressFromBech32(pool.Address)
		if err != nil {
			panic(err)
		}
		delAddr, err := sdk.AccAddressFromBech32(pool.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		var defiAddr sdk.ValAddress = nil
		if pool.DefiAddress != "" {
			defiAddr, err = sdk.ValAddressFromBech32(pool.DefiAddress)
			if err != nil {
				panic(err)
			}
		}
		keeper.CreatePool(ctx, poolAddr, delAddr, defiAddr, pool.Pledge, pool.UpdateTime)
	}

	for _, order := range data.Orders {
		ownerAddr, err := sdk.AccAddressFromBech32(order.OwnerAddress)
		if err != nil {
			panic(err)
		}
		poolAddr, err := sdk.AccAddressFromBech32(order.PoolAddress)
		if err != nil {
			panic(err)
		}
		keeper.PlaceOrder(ctx, poolAddr, ownerAddr, order.MyAsset, order.ExpectAsset, order.Price, order.Nonce)
	}

	for _, txPairStats := range data.TxPairsStats {
		keeper.SetTxPairStats(ctx, txPairStats)
	}

	return res
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, params, stats in
// the keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	var orders []types.Order
	
	keeper.IterateAllOrders(ctx, func(_ int64, order types.Order) (stop bool) {
		orders = append(orders, order)
		return false
	})

	var txPairsStats []types.TxPairStats
	keeper.IterateTxPairsStats(ctx, func(_ int64, txPairStats types.TxPairStats) (stop bool) {
		txPairsStats = append(txPairsStats, txPairStats)
		return false
	})

	return &types.GenesisState{
		Params:               keeper.GetParams(ctx),
		Pools:                keeper.GetAllPools(ctx),
		Orders:               orders,
		TxPairsStats:         txPairsStats,
		Exported:             true,
	}
}

// ValidateGenesis validates the provided staking genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data *types.GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}
	
	for _, pool := range data.Pools {
		if err := pool.Validate(); err != nil {
			return err
		}
	}

	for _, order := range data.Orders {
		if err := order.Validate(); err != nil {
			return err
		}
	}

	for _, txPairStats := range data.TxPairsStats {
		if err := txPairStats.Validate(); err != nil {
			return err
		}
	}

	return nil
}
