package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/gauss/gauss/v4/x/defi/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding staking type.
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.DefisKey):
			var defiA, defiB types.Defi

			cdc.MustUnmarshalBinaryBare(kvA.Value, &defiA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &defiB)

			return fmt.Sprintf("%v\n%v", defiA, defiB)

		case bytes.Equal(kvA.Key[:1], types.DelegationKey):
			var delegationA, delegationB types.Delegation

			cdc.MustUnmarshalBinaryBare(kvA.Value, &delegationA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &delegationB)

			return fmt.Sprintf("%v\n%v", delegationA, delegationB)

		case bytes.Equal(kvA.Key[:1], types.UnbondingDelegationKey),
			bytes.Equal(kvA.Key[:1], types.UnbondingDelegationByDefiIndexKey):
			var ubdA, ubdB types.UnbondingDelegation

			cdc.MustUnmarshalBinaryBare(kvA.Value, &ubdA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &ubdB)

			return fmt.Sprintf("%v\n%v", ubdA, ubdB)

		case bytes.Equal(kvA.Key[:1], types.DefiOutstandingRewardsPrefix):
			var rewardsA, rewardsB types.DefiOutstandingRewards
			cdc.MustUnmarshalBinaryBare(kvA.Value, &rewardsA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &rewardsB)
			return fmt.Sprintf("%v\n%v", rewardsA, rewardsB)

		case bytes.Equal(kvA.Key[:1], types.DelegatorWithdrawAddrPrefix):
			return fmt.Sprintf("%v\n%v", sdk.AccAddress(kvA.Value), sdk.AccAddress(kvB.Value))

		case bytes.Equal(kvA.Key[:1], types.DelegatorStartingInfoPrefix):
			var infoA, infoB types.DelegatorStartingInfo
			cdc.MustUnmarshalBinaryBare(kvA.Value, &infoA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &infoB)
			return fmt.Sprintf("%v\n%v", infoA, infoB)

		case bytes.Equal(kvA.Key[:1], types.DefiHistoricalRewardsPrefix):
			var rewardsA, rewardsB types.DefiHistoricalRewards
			cdc.MustUnmarshalBinaryBare(kvA.Value, &rewardsA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &rewardsB)
			return fmt.Sprintf("%v\n%v", rewardsA, rewardsB)

		case bytes.Equal(kvA.Key[:1], types.DefiCurrentRewardsPrefix):
			var rewardsA, rewardsB types.DefiCurrentRewards
			cdc.MustUnmarshalBinaryBare(kvA.Value, &rewardsA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &rewardsB)
			return fmt.Sprintf("%v\n%v", rewardsA, rewardsB)

		case bytes.Equal(kvA.Key[:1], types.DefiAccumulatedCommissionPrefix):
			var commissionA, commissionB types.DefiAccumulatedCommission
			cdc.MustUnmarshalBinaryBare(kvA.Value, &commissionA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &commissionB)
			return fmt.Sprintf("%v\n%v", commissionA, commissionB)


		default:
			panic(fmt.Sprintf("invalid staking key prefix %X", kvA.Key[:1]))
		}
	}
}
