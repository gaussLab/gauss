package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	
	"github.com/gauss/gauss/v4/x/token/types"
)

var _ Keeper = (*BaseKeeper)(nil)

// Keeper defines a module interface that facilitates the transfer of coins
// between accounts.
type Keeper interface {
	// basic
	SendKeeper

	// global
	InitGenesis(sdk.Context, *types.GenesisState)
	ExportGenesis(sdk.Context) *types.GenesisState

	// msgServer
	IssueToken(ctx sdk.Context, name string, symbol string, smallestUnit string, decimals uint32,initialSupply uint64,
		totalSupply uint64, mintable bool, unlocked bool, owner sdk.AccAddress) error 
	EditToken(ctx sdk.Context, symbol string, mintable bool, owner sdk.AccAddress) error 
	MintToken(ctx sdk.Context, symbol string, amount uint64, recipient sdk.AccAddress, owner sdk.AccAddress) error 
	MintTokenWithUnit(ctx sdk.Context, unit string, amount uint64, recipient string) error 
	BurnToken(ctx sdk.Context, symbol string, amount uint64, owner sdk.AccAddress) error 
	UnlockToken(ctx sdk.Context, symbol string, owner sdk.AccAddress) error
	TransferTokenOwner(ctx sdk.Context, symbol string, oldOwner sdk.AccAddress, newOwner sdk.AccAddress) error

	DeductIssueTokenFee(ctx sdk.Context, owner sdk.AccAddress, symbol string) error
	DeductMintTokenFee(ctx sdk.Context, owner sdk.AccAddress, symbol string) error
	GetIssueTokenFee(ctx sdk.Context, symbol string) (sdk.Coin, error)
	GetMintTokenFee(ctx sdk.Context, symbol string) (sdk.Coin, error)

	GetBlockedAddress()(map[string]bool)

	types.QueryServer
}

// BaseKeeper manages transfers between accounts. It implements the Keeper interface
type BaseKeeper struct {
	BaseSendKeeper

	storeKey		sdk.StoreKey
	cdc			codec.Marshaler
	bankKeeper		types.BankKeeper

	blockedAddress		map[string]bool

	feeCollectorName	string

	// params subspace
	paramSpace		paramstypes.Subspace
}

func NewKeeper(cdc codec.Marshaler, key sdk.StoreKey, paramSpace paramstypes.Subspace,
	bankKeeper types.BankKeeper, blockedAddress map[string]bool, feeCollectorName string) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return BaseKeeper{
		BaseSendKeeper:   NewBaseSendKeeper(cdc, key, bankKeeper, paramSpace),
		storeKey:         key,
		cdc:              cdc,
		paramSpace:       paramSpace,
		bankKeeper:       bankKeeper,
		blockedAddress:   blockedAddress,
		feeCollectorName: feeCollectorName,
	}
}

// IssueToken issues a new token
func (k BaseKeeper) IssueToken(
	ctx sdk.Context,
	name string,
	symbol string,
	smallestUnit string,
	decimals uint32,
	initialSupply uint64,
	totalSupply uint64,
	mintable bool,
	unlocked bool,
	owner sdk.AccAddress,
) error {

	// pools
	token := types.NewToken(
		name, symbol, smallestUnit, decimals, initialSupply,
		totalSupply, mintable, owner)
	if err := k.AddToken(ctx, token); err != nil {
		return err
	}

	// status
	if unlocked != k.IsUnlocked(ctx, smallestUnit) {
		k.unlockToken(ctx, smallestUnit, unlocked)
	}

	// circulation
	initialCoin := sdk.NewCoin(
		token.GetSmallestUnit(),
		sdk.NewIntFromUint64(token.GetInitialSupply()),
	)
	mintCoins := sdk.NewCoins(initialCoin)
	return k.mintCoinsToAccount(ctx, owner, mintCoins)
}

// EditToken edits the specified token
func (k BaseKeeper) EditToken(
	ctx sdk.Context,
	symbol string,
	mintable bool,
	owner sdk.AccAddress,
) error {
	token, err := k.getTokenBySymbol(ctx, symbol)
	if err != nil {
		return err
	}

	if owner.String() != token.GetOwnerString() {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, 
			"The address %s is not the owner of the token %s", owner, symbol)
	}
	
	token.Mintable = mintable
	k.storeToken(ctx, token)

	return nil
}

// MintToken mints the specified amount of token to the specified recipient
// NOTE: empty owner means that the external caller is responsible to manage the token authority
func (k BaseKeeper) MintToken(
	ctx sdk.Context,
	symbol string,
	amount uint64,
	recipient sdk.AccAddress,
	owner sdk.AccAddress,
) error {
	if amount == 0 {
		return nil
	}

	token, err := k.getTokenBySymbol(ctx, symbol)
	if err != nil {
		return err
	}

	if owner != nil {
		if owner.String() != token.GetOwnerString() {
			return sdkerrors.Wrapf(types.ErrInvalidOwner, 
				"the address %s is not the owner of the token %s", owner, symbol)
		}
	
		if recipient.Empty() {
			recipient = owner
		}
	}
	if recipient.Empty() {
		return types.ErrInvalidToAddress
	}

	coins, err := k.getMintCoins(ctx, token, amount)
	if err != nil {
		return err
	}
	return k.mintCoinsToAccount(ctx, recipient, coins)
}

func (k BaseKeeper) MintTokenWithUnit(
	ctx sdk.Context,
	unit string,
	amount uint64,
	recipient string,
) error {
	if amount == 0 {
		return nil
	}

	if recipient == "" {
		return types.ErrInvalidToAddress
	}

	token, err := k.GetTokenWithUnit(ctx, unit)
	if err != nil {
		return err
	}

	coins, err := k.getMintCoins(ctx, token, amount)
	if err != nil {
		return err
	}
	return k.mintCoinsToModule(ctx, recipient, coins)
}
func (k BaseKeeper) getMintCoins(
	ctx sdk.Context,
	token types.TokenI,
	amount uint64,
) (coins sdk.Coins, err error) {
	if !token.GetMintable() {
		return coins, sdkerrors.Wrapf(types.ErrNotMintable, "%s", token.GetSymbol())
	}

	totalSupply := k.getTokenSupply(ctx, token.GetSmallestUnit())
	burntCoin, found := k.GetBurntCoin(ctx, token.GetSmallestUnit())
	if !found {
		burntCoin = sdk.NewCoin(token.GetSmallestUnit(), sdk.NewInt(0))
	}
	mintableAmount := sdk.NewIntFromUint64(token.GetTotalSupply()).Sub(totalSupply).Sub(burntCoin.Amount)

	if amount > mintableAmount.Uint64() {
		return coins, sdkerrors.Wrapf(
			types.ErrInvalidAmount,
			"the amount exceeds the mintable token amount; expected (0, %d], got %d",
			mintableAmount, amount,
		)
	}

	mintCoin := sdk.NewCoin(token.GetSmallestUnit(), sdk.NewIntFromUint64(amount))
	return sdk.NewCoins(mintCoin), nil
}


// BurnToken burns the specified amount of token
func (k BaseKeeper) BurnToken(
	ctx sdk.Context,
	symbol string,
	amount uint64,
	owner sdk.AccAddress,
) error {
	token, err := k.getTokenBySymbol(ctx, symbol)
	if err != nil {
		return err
	}
	if token.GetOwnerString() != owner.String() {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "account %s does not have permissions to burn tokens", token.GetOwnerString())
	}

	burnedCoin := sdk.NewCoin(token.GetSmallestUnit(), sdk.NewIntFromUint64(amount))
	burnedCoins := sdk.NewCoins(burnedCoin)
	if err := k.burnCoins(ctx, owner, burnedCoins); err != nil {
		return err
	}
	
	k.AddBurnedCoin(ctx, burnedCoin)
	return nil
}

// UnlockToken unlock the specialied token
func (k BaseKeeper) UnlockToken(ctx sdk.Context, symbol string, owner sdk.AccAddress) error {
	token, err := k.getTokenBySymbol(ctx, symbol)
	if err != nil {
		return err
	}
	
	if token.GetOwnerString() != owner.String() {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "%s is not the owner of the token[%s]", 
			owner.String(), symbol)
	}

	if k.IsUnlocked(ctx, token.GetSmallestUnit()) {
		return sdkerrors.Wrapf(types.ErrUnlockedToken, "the token[%s] has been unlocked.",
			token.GetSmallestUnit())
	}

	k.unlockToken(ctx, token.GetSmallestUnit(), true)
	return nil
}

// TransferTokenOwner transfers the owner of the specified token to a new one
func (k BaseKeeper) TransferTokenOwner(
	ctx sdk.Context,
	symbol string,
	oldOwner sdk.AccAddress,
	newOwner sdk.AccAddress,
) error {
	token, err := k.getTokenBySymbol(ctx, symbol)
	if err != nil {
		return err
	}

	if oldOwner.String() != token.GetOwnerString() {
		return sdkerrors.Wrapf(types.ErrInvalidOwner, "%s is not the owner of the token[%s]", 
			oldOwner, symbol)
	}

	// update token
	token.Owner = newOwner.String()
	k.storeToken(ctx, token)
	k.resetTokenOwner(ctx, token.GetSymbol(), oldOwner, newOwner)

	return nil
}

func (k BaseKeeper) GetBlockedAddress()(map[string]bool){
	return k.blockedAddress;
}
