package keeper

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/gauss/gauss/v4/x/defi/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// Defis queries all defis that match the given status
func (k Querier) Defis(c context.Context, req *types.QueryDefisRequest) (*types.QueryDefisResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// validate the provided status, return all the defis if the status is empty
	if req.Status != "" && !(req.Status == types.Bonded.String() || req.Status == types.Unbonded.String() || req.Status == types.Unbonding.String()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid defi status %s", req.Status)
	}

	var defis types.Defis
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	defiStore := prefix.NewStore(store, types.DefisKey)

	pageRes, err := query.FilteredPaginate(defiStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		defi, err := types.UnmarshalDefi(k.cdc, value)
		if err != nil {
			return false, err
		}

		if req.Status != "" && !strings.EqualFold(defi.GetStatus().String(), req.Status) {
			return false, nil
		}

		if accumulate {
			defis = append(defis, defi)
		}

		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDefisResponse{Defis: defis, Pagination: pageRes}, nil
}

// Defi queries defi info for given defi addr
func (k Querier) Defi(c context.Context, req *types.QueryDefiRequest) (*types.QueryDefiResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DefiAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "defi address cannot be empty")
	}

	defiAddr, err := sdk.ValAddressFromBech32(req.DefiAddress)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	defi, found := k.GetDefi(ctx, defiAddr)
	if !found {
		return nil, status.Errorf(codes.NotFound, "defi %s not found", req.DefiAddress)
	}

	return &types.QueryDefiResponse{Defi: defi}, nil
}

// DefiDelegations queries delegate info for given defi
func (k Querier) DefiDelegations(c context.Context, req *types.QueryDefiDelegationsRequest) (*types.QueryDefiDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DefiAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "defi address cannot be empty")
	}
	var delegations []types.Delegation
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	defiStore := prefix.NewStore(store, types.DelegationKey)
	pageRes, err := query.FilteredPaginate(defiStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return false, err
		}

		defiAddr, err := sdk.ValAddressFromBech32(req.DefiAddress)
		if err != nil {
			return false, err
		}

		if !delegation.GetDefiAddr().Equals(defiAddr) {
			return false, nil
		}

		if accumulate {
			delegations = append(delegations, delegation)
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	delResponses, err := DelegationsToDelegationResponses(ctx, k.Keeper, delegations)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDefiDelegationsResponse{
		DelegationResponses: delResponses, Pagination: pageRes}, nil
}

// DefiUnbondingDelegations queries unbonding delegations of a defi
func (k Querier) DefiUnbondingDelegations(c context.Context, req *types.QueryDefiUnbondingDelegationsRequest) (*types.QueryDefiUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DefiAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "defi address cannot be empty")
	}
	var ubds types.UnbondingDelegations
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)

	defiAddr, err := sdk.ValAddressFromBech32(req.DefiAddress)
	if err != nil {
		return nil, err
	}

	srcDefiPrefix := types.GetUBDsByDefiIndexKey(defiAddr)
	ubdStore := prefix.NewStore(store, srcDefiPrefix)
	pageRes, err := query.Paginate(ubdStore, req.Pagination, func(key []byte, value []byte) error {
		storeKey := types.GetUBDKeyFromDefiIndexKey(append(srcDefiPrefix, key...))
		storeValue := store.Get(storeKey)

		ubd, err := types.UnmarshalUBD(k.cdc, storeValue)
		if err != nil {
			return err
		}
		ubds = append(ubds, ubd)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDefiUnbondingDelegationsResponse{
		UnbondingResponses: ubds,
		Pagination:         pageRes,
	}, nil
}

// Delegation queries delegate info for given defil delegator pair
func (k Querier) DefiDelegation(c context.Context, req *types.QueryDelegationRequest) (*types.QueryDelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.DefiAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "defi address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	defiAddr, err := sdk.ValAddressFromBech32(req.DefiAddress)
	if err != nil {
		return nil, err
	}

	delegation, found := k.GetDelegation(ctx, delAddr, defiAddr)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"delegation with delegator %s not found for defi %s",
			req.DelegatorAddress, req.DefiAddress)
	}

	delResponse, err := DelegationToDelegationResponse(ctx, k.Keeper, delegation)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegationResponse{DelegationResponse: &delResponse}, nil
}

// UnbondingDelegation queries unbonding info for give defi delegator pair
func (k Querier) DefiUnbondingDelegation(c context.Context, req *types.QueryUnbondingDelegationRequest) (*types.QueryUnbondingDelegationResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Errorf(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.DefiAddress == "" {
		return nil, status.Errorf(codes.InvalidArgument, "defi address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	defiAddr, err := sdk.ValAddressFromBech32(req.DefiAddress)
	if err != nil {
		return nil, err
	}

	unbond, found := k.GetUnbondingDelegation(ctx, delAddr, defiAddr)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"unbonding delegation with delegator %s not found for defi %s",
			req.DelegatorAddress, req.DefiAddress)
	}

	return &types.QueryUnbondingDelegationResponse{Unbond: unbond}, nil
}

// DelegatorDelegations queries all delegations of a give delegator address
func (k Querier) DefiDelegatorDelegations(c context.Context, req *types.QueryDelegatorDelegationsRequest) (*types.QueryDelegatorDelegationsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	var delegations types.Delegations
	ctx := sdk.UnwrapSDKContext(c)

	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	store := ctx.KVStore(k.storeKey)
	delStore := prefix.NewStore(store, types.GetDelegationsKey(delAddr))
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}
		delegations = append(delegations, delegation)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if delegations == nil {
		return nil, status.Errorf(
			codes.NotFound,
			"unable to find delegations for address %s", req.DelegatorAddress)
	}
	delegationResps, err := DelegationsToDelegationResponses(ctx, k.Keeper, delegations)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorDelegationsResponse{DelegationResponses: delegationResps, Pagination: pageRes}, nil

}

// DelegatorDefi queries defi info for given delegator defi pair
func (k Querier) DelegatorDefi(c context.Context, req *types.QueryDelegatorDefiRequest) (*types.QueryDelegatorDefiResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.DefiAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "defi address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	defiAddr, err := sdk.ValAddressFromBech32(req.DefiAddress)
	if err != nil {
		return nil, err
	}

	defi, err := k.GetDelegatorDefi(ctx, delAddr, defiAddr)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorDefiResponse{Defi: defi}, nil
}

func (k Querier) DelegatorDefis(c context.Context, req *types.QueryDelegatorDefisRequest) (*types.QueryDelegatorDefisResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	var defis types.Defis
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	delStore := prefix.NewStore(store, types.GetDelegationsKey(delAddr))
	pageRes, err := query.Paginate(delStore, req.Pagination, func(key []byte, value []byte) error {
		delegation, err := types.UnmarshalDelegation(k.cdc, value)
		if err != nil {
			return err
		}

		defi, found := k.GetDefi(ctx, delegation.GetDefiAddr())
		if !found {
			return types.ErrNoDefiFound
		}

		defis = append(defis, defi)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorDefisResponse{Defis: defis, Pagination: pageRes}, nil
}

// DelegatorUnbondingDelegations queries all unbonding delegations of a given delegator address
func (k Querier) DefiDelegatorUnbondingDelegations(c context.Context, req *types.QueryDelegatorUnbondingDelegationsRequest) (*types.QueryDelegatorUnbondingDelegationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	var unbondingDelegations types.UnbondingDelegations
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	delAddr, err := sdk.AccAddressFromBech32(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	unbStore := prefix.NewStore(store, types.GetUBDsKey(delAddr))
	pageRes, err := query.Paginate(unbStore, req.Pagination, func(key []byte, value []byte) error {
		unbond, err := types.UnmarshalUBD(k.cdc, value)
		if err != nil {
			return err
		}
		unbondingDelegations = append(unbondingDelegations, unbond)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDelegatorUnbondingDelegationsResponse{
		UnbondingResponses: unbondingDelegations, Pagination: pageRes}, nil
}

// HistoricalInfo queries the historical info for given height
func (k Querier) DefiHistoricalInfo(c context.Context, req *types.QueryHistoricalInfoRequest) (*types.QueryHistoricalInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Height < 0 {
		return nil, status.Error(codes.InvalidArgument, "height cannot be negative")
	}
	ctx := sdk.UnwrapSDKContext(c)
	hi, found := k.GetHistoricalInfo(ctx, req.Height)
	if !found {
		return nil, status.Errorf(codes.NotFound, "historical info for height %d not found", req.Height)
	}

	return &types.QueryHistoricalInfoResponse{Hist: &hi}, nil
}

// Pool queries the pool info
func (k Querier) DefiPool(c context.Context, _ *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	bondDenom := k.BondDenom(ctx)
	bondedPool := k.GetBondedPool(ctx)
	notBondedPool := k.GetNotBondedPool(ctx)

	pool := types.NewPool(
		k.bankKeeper.GetBalance(ctx, notBondedPool.GetAddress(), bondDenom).Amount,
		k.bankKeeper.GetBalance(ctx, bondedPool.GetAddress(), bondDenom).Amount,
	)

	return &types.QueryPoolResponse{Pool: pool}, nil
}

// Params queries the staking parameters
func (k Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// Rewards
// DefiOutstandingRewards queries rewards of a defi address
func (k Querier) DefiOutstandingRewards(c context.Context, req *types.QueryDefiOutstandingRewardsRequest) (*types.QueryDefiOutstandingRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.DefiAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty defi address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	defiAdr, err := sdk.ValAddressFromBech32(req.DefiAddress)
	if err != nil {
		return nil, err
	}
	rewards := k.GetDefiOutstandingRewards(ctx, defiAdr)

	return &types.QueryDefiOutstandingRewardsResponse{Rewards: rewards}, nil
}

// DefiCommission queries accumulated commission for a defi
func (k Querier) DefiCommission(c context.Context, req *types.QueryDefiCommissionRequest) (*types.QueryDefiCommissionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.DefiAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty defi address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	defiAdr, err := sdk.ValAddressFromBech32(req.DefiAddress)
	if err != nil {
		return nil, err
	}
	commission := k.GetDefiAccumulatedCommission(ctx, defiAdr)

	return &types.QueryDefiCommissionResponse{Commission: commission}, nil
}

// DelegationRewards the total rewards accrued by a delegation
func (k Querier) DefiDelegationRewards(c context.Context, req *types.QueryDelegationRewardsRequest) (*types.QueryDelegationRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delegator address")
	}

	if req.DefiAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty defi address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	defiAdr, err := sdk.ValAddressFromBech32(req.DefiAddress)
	if err != nil {
		return nil, err
	}

	defi := k.Keeper.Defi(ctx, defiAdr)
	if defi == nil {
		return nil, sdkerrors.Wrap(types.ErrNoDefiFound, req.DefiAddress)
	}

	delAdr, err := sdk.AccAddressFromBech32(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	del := k.Keeper.Delegation(ctx, delAdr, defiAdr)
	if del == nil {
		return nil, types.ErrBadDelegatorAddr
	}

	endingPeriod := k.IncrementDefiPeriod(ctx, defi)
	rewards := k.CalculateDelegationRewards(ctx, defi, del, endingPeriod)

	return &types.QueryDelegationRewardsResponse{Rewards: rewards}, nil
}

// DelegationTotalRewards the total rewards accrued by a each defi
func (k Querier) DefiDelegationTotalRewards(c context.Context, req *types.QueryDelegationTotalRewardsRequest) (*types.QueryDelegationTotalRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delegator address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	total := sdk.DecCoins{}
	var delRewards []types.DelegationDelegatorReward

	delAdr, err := sdk.AccAddressFromBech32(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	k.Keeper.IterateDelegations(
		ctx, delAdr,
		func(_ int64, del types.DelegationI) (stop bool) {
			defiAddr := del.GetDefiAddr()
			defi := k.Keeper.Defi(ctx, defiAddr)
			endingPeriod := k.IncrementDefiPeriod(ctx, defi)
			delReward := k.CalculateDelegationRewards(ctx, defi, del, endingPeriod)

			delRewards = append(delRewards, types.NewDelegationDelegatorReward(defiAddr, delReward))
			total = total.Add(delReward...)
			return false
		},
	)

	return &types.QueryDelegationTotalRewardsResponse{Rewards: delRewards, Total: total}, nil
}

// DelegatorDefisEx queries the defis list of a delegator
func (k Querier) DelegatorDefisEx(c context.Context, req *types.QueryDelegatorDefisExRequest) (*types.QueryDelegatorDefisExResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delegator address")
	}

	ctx := sdk.UnwrapSDKContext(c)
	delAdr, err := sdk.AccAddressFromBech32(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	var defis []string

	k.Keeper.IterateDelegations(
		ctx, delAdr,
		func(_ int64, del types.DelegationI) (stop bool) {
			defis = append(defis, del.GetDefiAddr().String())
			return false
		},
	)

	return &types.QueryDelegatorDefisExResponse{Defis: defis}, nil
}

// DelegatorWithdrawAddress queries Query/delegatorWithdrawAddress
func (k Querier) DefiDelegatorWithdrawAddress(c context.Context, req *types.QueryDelegatorWithdrawAddressRequest) (*types.QueryDelegatorWithdrawAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.DelegatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty delegator address")
	}
	delAdr, err := sdk.AccAddressFromBech32(req.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, delAdr)

	return &types.QueryDelegatorWithdrawAddressResponse{WithdrawAddress: withdrawAddr.String()}, nil
}

// CommunityPool queries the community pool coins
func (k Querier) DefiCommunityPool(c context.Context, req *types.QueryCommunityPoolRequest) (*types.QueryCommunityPoolResponse, error) {
        ctx := sdk.UnwrapSDKContext(c)
        pool := k.GetFeePoolCommunityCoins(ctx)

        return &types.QueryCommunityPoolResponse{Pool: pool}, nil
}
