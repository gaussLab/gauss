package keeper

import (
	"context"
	"fmt"

	"github.com/gogo/protobuf/proto"
	gogotypes "github.com/gogo/protobuf/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
//	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/gauss/gauss/v4/x/token/types"
)

var _ types.QueryServer = BaseKeeper{}

// Params return the all the parameter in tonken module
func (k BaseKeeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (k BaseKeeper) Token(c context.Context, req *types.QueryTokenRequest) (*types.QueryTokenResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if err := types.ValidateSymbol(req.Symbol); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)

	token, err := k.GetToken(ctx, req.Symbol)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Token %s do not found", req.Symbol)
	}

	msg, ok := token.(proto.Message)
	if !ok {
		return nil, status.Errorf(codes.Internal, "can't protomarshal %T", token)
	}
	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryTokenResponse{
		Token: any,
		Unlocked: k.IsUnlocked(ctx, token.GetSmallestUnit()),
		}, nil
}

func (k BaseKeeper) Tokens(c context.Context, req *types.QueryTokensRequest) (*types.QueryTokensResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var owner sdk.AccAddress
	var err error
	if req.Owner != ""  {
		owner, err = sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("invalid owner address: %s", err.Error()))
		}
	}

	var tokens []types.TokenI
	var pageRes *query.PageResponse

	store := ctx.KVStore(k.storeKey)
	if owner == nil {
		tokenStore := prefix.NewStore(store, types.SymbolPrefix)
		pageRes, err = query.Paginate(tokenStore, req.Pagination, func(key []byte, value []byte) error {
			var token types.Token

			k.cdc.MustUnmarshalBinaryBare(value, &token)
			tokens = append(tokens, &token)

			return nil
		})
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
		}
	} else {
		tokenStore := prefix.NewStore(store, types.GetOwnerSymbolKey(owner, ""))
		pageRes, err = query.Paginate(tokenStore, req.Pagination, func(key []byte, value []byte) error {
			var symbol gogotypes.StringValue

			k.cdc.MustUnmarshalBinaryBare(value, &symbol)
			token, err := k.GetToken(ctx, symbol.Value)
			if err == nil {
				tokens = append(tokens, token)
			}

			return nil
		})
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "paginate: %v", err)
		}
	}

	result := make([]*codectypes.Any, len(tokens))
	for i, token := range tokens {
		msg, ok := token.(proto.Message)
		if !ok {
			return nil, status.Errorf(codes.Internal, "%T does not implement proto.Message", token)
		}

		if result[i], err = codectypes.NewAnyWithValue(msg); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &types.QueryTokensResponse{Tokens: result, Pagination: pageRes}, nil
}

func (k BaseKeeper) Burntoken(c context.Context, req *types.QueryBurntokenRequest) (*types.QueryBurntokenResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if err := types.ValidateSymbol(req.Symbol); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(c)

	token, err := k.GetToken(ctx, req.Symbol)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Token %s do not found", req.Symbol)	
	}

	burntCoin, found := k.GetBurntCoin(ctx, token.GetSmallestUnit())
	if !found {
		return &types.QueryBurntokenResponse{
			BurnedCoin: sdk.Coin{Denom:token.GetSmallestUnit(), Amount:sdk.ZeroInt()},
		}, nil
		// return nil, sdkerrors.Wrapf(err, "failed to get burned coin of the token[%s].", req.Symbol)
	}

	return &types.QueryBurntokenResponse{
		Exist:      k.HasToken(ctx, req.Symbol),
		BurnedCoin: burntCoin,
	}, nil
}

func (k BaseKeeper) Fees(c context.Context, req *types.QueryFeesRequest) (*types.QueryFeesResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := types.ValidateSymbol(req.Symbol); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	issueFee, err := k.GetIssueTokenFee(ctx, req.Symbol)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	mintFee, err := k.GetMintTokenFee(ctx, req.Symbol)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryFeesResponse{
		Exist:    k.HasToken(ctx, req.Symbol),
		IssueFee: issueFee,
		MintFee:  mintFee,
	}, nil
}
