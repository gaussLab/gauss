package orderbook

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	// govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	// ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"

	"github.com/gauss/gauss/v4/x/orderbook/keeper"
)


// ValidateOrderbookDecorator is responsible for restricting the token participation of the swap prefix
type ValidateOrderbookDecorator struct {
	keeper keeper.Keeper
}

// NewValidateOrderbookDecorator returns an instance of ValidateOrderbookDecorator
func NewValidateOrderbookDecorator(tk keeper.Keeper) ValidateOrderbookDecorator {
	return ValidateOrderbookDecorator {
		keeper: tk,
	}
}

// AnteHandle checks the transaction
func (vtd ValidateOrderbookDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, 
	simulate bool, next sdk.AnteHandler) (sdk.Context, error) {

/*
	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *types.MsgCreatePool:
		case *ibctransfertypes.MsgTransfer:
		case *govtypes.MsgSubmitProposal:
		case *govtypes.MsgDeposit:
		default:
			break
		}
	}
*/

	return next(ctx, tx, simulate)
}
