package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gauss/gauss/v4/x/defi/types"
)

// AllocateTokens handles distribution of the collected fees
func (k Keeper) MintTokens(
	ctx sdk.Context, defiAddr sdk.ValAddress, marketAddr string, marketRate sdk.Dec, toModule bool,
) {
	// logger := k.Logger(ctx)

	// all rewards
	mintedTokens := k.mintTokens(ctx, defiAddr)
	if mintedTokens.IsZero() {
		return 
	}

	// market rate
	marketRateL := k.MarketRate(ctx)
	if marketRate.LT(marketRateL) {
		marketRateL = marketRate
	}
		
	// community tax
	communityTax := k.CommunityTax(ctx)

	// delegator rate
	voteMultiplier := sdk.OneDec().Sub(marketRateL).Sub(communityTax)

	// market rewards
	marketTokens := mintedTokens.MulDec(marketRateL)
	k.allocateTokensToMarket(ctx, marketAddr, marketTokens, toModule)

	// delegator rewards
	defi := k.Defi(ctx, defiAddr)
	rewards := mintedTokens.MulDec(voteMultiplier)
	k.allocateTokensToDefi(ctx, defi, rewards)
	remaining := mintedTokens.Sub(marketTokens).Sub(rewards)

	// allocate community funding
	// temporary workaround to keep CanWithdrawInvariant happy
	// general discussions here: https://github.com/cosmos/cosmos-sdk/issues/2906#issuecomment-441867634
	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
	k.SetFeePool(ctx, feePool)
}

func (k Keeper) mintTokens(
	ctx sdk.Context, defiAddr sdk.ValAddress,
) sdk.DecCoins {
	mintedCoin := k.transactionProvision(ctx)

	if !mintedCoin.IsZero() {
		if err := k.mintCoin(ctx, mintedCoin) ; err != nil {
			return sdk.NewDecCoins(sdk.NewDecCoin(mintedCoin.Denom, sdk.NewInt(0)))
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeMint,
				sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.String()),
				sdk.NewAttribute(types.AttributeKeyDefi, defiAddr.String()),
                	),
        	)
	}

	return sdk.NewDecCoinsFromCoins(sdk.NewCoins(mintedCoin)...)
}

// allocateTokensToMarket implements an alias call to the underlying supply keeper's
func (k Keeper) allocateTokensToMarket(ctx sdk.Context, marketAddrString string, tokens sdk.DecCoins, toModule bool) error {
	coins, _ := tokens.TruncateDecimal()
	if toModule {
		err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, marketAddrString, coins)
		if err != nil {
			return err
		}
	} else {
		marketAddr, err := sdk.AccAddressFromBech32(marketAddrString)
		if err != nil {
			return err
		}
		
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, marketAddr, coins)
		if err != nil {
			return err
		}
	}

	// update account rewards
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, coins.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, marketAddrString),
		),
	)

	return nil
}


// allocateTokensToDefi allocate tokens to a particular defi, splitting according to commission
func (k Keeper) allocateTokensToDefi(ctx sdk.Context, defi types.DefiI, tokens sdk.DecCoins) {
	// split tokens between defi and delegators according to commission
	commission := tokens.MulDec(k.CommissionRate(ctx))
	shared := tokens.Sub(commission)

	// update current commission
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommission,
			sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
			sdk.NewAttribute(types.AttributeKeyDefi, defi.GetOperator().String()),
		),
	)
	currentCommission := k.GetDefiAccumulatedCommission(ctx, defi.GetOperator())
	currentCommission.Commission = currentCommission.Commission.Add(commission...)
	k.SetDefiAccumulatedCommission(ctx, defi.GetOperator(), currentCommission)

	// update current rewards
	currentRewards := k.GetDefiCurrentRewards(ctx, defi.GetOperator())
	currentRewards.Rewards = currentRewards.Rewards.Add(shared...)
	k.SetDefiCurrentRewards(ctx, defi.GetOperator(), currentRewards)

	// update outstanding rewards
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, tokens.String()),
			sdk.NewAttribute(types.AttributeKeyDefi, defi.GetOperator().String()),
		),
	)
	outstanding := k.GetDefiOutstandingRewards(ctx, defi.GetOperator())
	outstanding.Rewards = outstanding.Rewards.Add(tokens...)
	k.SetDefiOutstandingRewards(ctx, defi.GetOperator(), outstanding)
}

func (k Keeper) transactionProvision(ctx sdk.Context) sdk.Coin {
	// if k.DefiTokenSupply(ctx).GT(defiTotalSupply){
	//	return sdk.NewCoin(k.DefiBondDenom(ctx), sdk.NewInt(0))
	//}
	
	return k.MintInflation(ctx)
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in per Transaction.
func (k Keeper) mintCoin(ctx sdk.Context, newCoin sdk.Coin) error {
	if newCoin.IsZero() {
		// skip as no coins need to be minted
		return nil
	}

	return k.tokenKeeper.MintTokenWithUnit(ctx, newCoin.Denom, newCoin.Amount.Uint64(), types.ModuleName)
}
