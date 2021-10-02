package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gauss/gauss/v4/x/defi/types"
)

//______________________________________________________
// util

func DelegationToDelegationResponse(ctx sdk.Context, k Keeper, del types.Delegation) (types.DelegationResponse, error) {
	defi, found := k.GetDefi(ctx, del.GetDefiAddr())
	if !found {
		return types.DelegationResponse{}, types.ErrNoDefiFound
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(del.DelegatorAddress)
	if err != nil {
		return types.DelegationResponse{}, err
	}

	return types.NewDelegationResp(
		delegatorAddress,
		del.GetDefiAddr(),
		del.Shares,
		sdk.NewCoin(k.BondDenom(ctx), defi.TokensFromShares(del.Shares).TruncateInt()),
	), nil
}

func DelegationsToDelegationResponses(
	ctx sdk.Context, k Keeper, delegations types.Delegations,
) (types.DelegationResponses, error) {
	resp := make(types.DelegationResponses, len(delegations))

	for i, del := range delegations {
		delResp, err := DelegationToDelegationResponse(ctx, k, del)
		if err != nil {
			return nil, err
		}

		resp[i] = delResp
	}

	return resp, nil
}
