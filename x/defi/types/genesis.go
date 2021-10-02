package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// NewGenesisState creates a new GenesisState instanc e
func NewGenesisState(params Params, fp FeePool, defis []Defi, delegations []Delegation,
	dwis []DelegatorWithdrawInfo, r []DefiOutstandingRewardsRecord, 
	acc []DefiAccumulatedCommissionRecord, historical []DefiHistoricalRewardsRecord, 
	cur []DefiCurrentRewardsRecord, dels []DelegatorStartingInfoRecord,
) *GenesisState {
	return &GenesisState{
		Params:                          params,
		FeePool:                         fp,
		Defis:	                         defis,
		Delegations: 		         delegations,
		DelegatorWithdrawInfos:          dwis,
		OutstandingRewards:              r,
		DefiAccumulatedCommissions:      acc,
		DefiHistoricalRewards:           historical,
		DefiCurrentRewards:              cur,
		DelegatorStartingInfos:          dels,
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		FeePool:                         InitialFeePool(),
		Params: 			 DefaultParams(),
		DelegatorWithdrawInfos:          []DelegatorWithdrawInfo{},
		OutstandingRewards:              []DefiOutstandingRewardsRecord{},
		DefiAccumulatedCommissions:      []DefiAccumulatedCommissionRecord{},
		DefiHistoricalRewards:           []DefiHistoricalRewardsRecord{},
		DefiCurrentRewards:              []DefiCurrentRewardsRecord{},
		DelegatorStartingInfos:          []DelegatorStartingInfoRecord{},
	}
}

// GetGenesisStateFromAppState returns x/defi GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.JSONMarshaler, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (g GenesisState) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	for i := range g.Defis {
		if err := g.Defis[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	return nil
}
