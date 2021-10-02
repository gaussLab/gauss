package keeper // noalias

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gauss/gauss/v4/x/defi/types"
)

// update defi for testing
func TestingUpdateDefi(keeper Keeper, ctx sdk.Context, defi types.Defi, apply bool) types.Defi {
	keeper.SetDefi(ctx, defi)

	if !apply {
		ctx, _ = ctx.CacheContext()
	}
	_, err := keeper.ApplyAndReturnDefiSetUpdates(ctx)
	if err != nil {
		panic(err)
	}

	defi, found := keeper.GetDefi(ctx, defi.GetOperator())
	if !found {
		panic("defi expected but not found")
	}

	return defi
}

// RandomDefi returns a random defi given access to the keeper and ctx
func RandomDefi(r *rand.Rand, keeper Keeper, ctx sdk.Context) (val types.Defi, ok bool) {
	defis := keeper.GetAllDefis(ctx)
	if len(defis) == 0 {
		return types.Defi{}, false
	}

	i := r.Intn(len(defis))

	return defis[i], true
}
