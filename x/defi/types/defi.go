package types

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// TODO: Why can't we just have one string description which can be JSON by convention
	MaxMonikerLength         = 70
	MaxIdentityLength        = 3000
	MaxWebsiteLength         = 140
	MaxSecurityContactLength = 140
	MaxDetailsLength         = 280
)

var (
        BondStatusUnspecified = BondStatus_name[int32(Unspecified)]
	BondStatusUnbonded    = BondStatus_name[int32(Unbonded)]
	BondStatusUnbonding   = BondStatus_name[int32(Unbonding)]
	BondStatusBonded      = BondStatus_name[int32(Bonded)]
)

var _ DefiI = Defi{}

// NewDefi constructs a new Defi
//nolint:interfacer
func NewDefi(operator sdk.ValAddress, description Description) (Defi, error) {
	return Defi{
		OperatorAddress:   operator.String(),
		Status:            Unbonded,
		Tokens:            sdk.ZeroInt(),
		DelegatorShares:   sdk.ZeroDec(),
		Description:       description,
		UnbondingHeight:   int64(0),
		UnbondingTime:     time.Unix(0, 0).UTC(),
		MinSelfDelegation: sdk.OneInt(),
	}, nil
}

// String implements the Stringer interface for a Defi object.
func (v Defi) String() string {
	out, _ := yaml.Marshal(v)
	return string(out)
}

// Defis is a collection of Defi
type Defis []Defi

func (v Defis) String() (out string) {
	for _, val := range v {
		out += val.String() + "\n"
	}

	return strings.TrimSpace(out)
}

// ToSDKDefis -  convenience function convert []Defi to []sdk.DefiI
func (v Defis) ToSDKDefis() (defis []DefiI) {
	for _, val := range v {
		defis = append(defis, val)
	}

	return defis
}

// Sort Defis sorts defi array in ascending operator address order
func (v Defis) Sort() {
	sort.Sort(v)
}

// Implements sort interface
func (v Defis) Len() int {
	return len(v)
}

// Implements sort interface
func (v Defis) Less(i, j int) bool {
	return bytes.Compare(v[i].GetOperator().Bytes(), v[j].GetOperator().Bytes()) == -1
}

// Implements sort interface
func (v Defis) Swap(i, j int) {
	it := v[i]
	v[i] = v[j]
	v[j] = it
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (v Defis) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	for i := range v {
		if err := v[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	return nil
}

// return the delegation
func MustMarshalDefi(cdc codec.BinaryMarshaler, defi *Defi) []byte {
        return cdc.MustMarshalBinaryBare(defi)
}

// unmarshal a delegation from a store value
func MustUnmarshalDefi(cdc codec.BinaryMarshaler, value []byte) Defi {
        defi, err := UnmarshalDefi(cdc, value)
        if err != nil {
                panic(err)
        }

        return defi
}

// unmarshal a delegation from a store value
func UnmarshalDefi(cdc codec.BinaryMarshaler, value []byte) (v Defi, err error) {
        err = cdc.UnmarshalBinaryBare(value, &v)
        return v, err
}

// IsBonded checks if the defi status equals Bonded
func (v Defi) IsBonded() bool {
	return v.GetStatus() == Bonded
}

// IsUnbonded checks if the defi status equals Unbonded
func (v Defi) IsUnbonded() bool {
	return v.GetStatus() == Unbonded
}

// IsUnbonding checks if the defi status equals Unbonding
func (v Defi) IsUnbonding() bool {
	return v.GetStatus() == Unbonding
}

// constant used in flags to indicate that description field should not be updated
const DoNotModifyDesc = "[do-not-modify]"

func NewDescription(moniker, identity, website, securityContact, details string) Description {
	return Description{
		Moniker:         moniker,
		Identity:        identity,
		Website:         website,
		SecurityContact: securityContact,
		Details:         details,
	}
}

// String implements the Stringer interface for a Description object.
func (d Description) String() string {
	out, _ := yaml.Marshal(d)
	return string(out)
}

// UpdateDescription updates the fields of a given description. An error is
// returned if the resulting description contains an invalid length.
func (d Description) UpdateDescription(d2 Description) (Description, error) {
	if d2.Moniker == DoNotModifyDesc {
		d2.Moniker = d.Moniker
	}

	if d2.Identity == DoNotModifyDesc {
		d2.Identity = d.Identity
	}

	if d2.Website == DoNotModifyDesc {
		d2.Website = d.Website
	}

	if d2.SecurityContact == DoNotModifyDesc {
		d2.SecurityContact = d.SecurityContact
	}

	if d2.Details == DoNotModifyDesc {
		d2.Details = d.Details
	}

	return NewDescription(
		d2.Moniker,
		d2.Identity,
		d2.Website,
		d2.SecurityContact,
		d2.Details,
	).EnsureLength()
}

// EnsureLength ensures the length of a defi's description.
func (d Description) EnsureLength() (Description, error) {
	if len(d.Moniker) > MaxMonikerLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid moniker length; got: %d, max: %d", len(d.Moniker), MaxMonikerLength)
	}

	if len(d.Identity) > MaxIdentityLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid identity length; got: %d, max: %d", len(d.Identity), MaxIdentityLength)
	}

	if len(d.Website) > MaxWebsiteLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid website length; got: %d, max: %d", len(d.Website), MaxWebsiteLength)
	}

	if len(d.SecurityContact) > MaxSecurityContactLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid security contact length; got: %d, max: %d", len(d.SecurityContact), MaxSecurityContactLength)
	}

	if len(d.Details) > MaxDetailsLength {
		return d, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid details length; got: %d, max: %d", len(d.Details), MaxDetailsLength)
	}

	return d, nil
}

// In some situations, the exchange rate becomes invalid, e.g. if
// Defi loses all tokens due to slashing. In this case,
// make all future delegations invalid.
func (v Defi) InvalidExRate() bool {
	return v.Tokens.IsZero() && v.DelegatorShares.IsPositive()
}

// calculate the token worth of provided shares
func (v Defi) TokensFromShares(shares sdk.Dec) sdk.Dec {
	return (shares.MulInt(v.Tokens)).Quo(v.DelegatorShares)
}

// calculate the token worth of provided shares, truncated
func (v Defi) TokensFromSharesTruncated(shares sdk.Dec) sdk.Dec {
	return (shares.MulInt(v.Tokens)).QuoTruncate(v.DelegatorShares)
}

// TokensFromSharesRoundUp returns the token worth of provided shares, rounded
// up.
func (v Defi) TokensFromSharesRoundUp(shares sdk.Dec) sdk.Dec {
	return (shares.MulInt(v.Tokens)).QuoRoundUp(v.DelegatorShares)
}

// SharesFromTokens returns the shares of a delegation given a bond amount. It
// returns an error if the defi has no tokens.
func (v Defi) SharesFromTokens(amt sdk.Int) (sdk.Dec, error) {
	if v.Tokens.IsZero() {
		return sdk.ZeroDec(), ErrInsufficientShares
	}

	return v.GetDelegatorShares().MulInt(amt).QuoInt(v.GetTokens()), nil
}

// SharesFromTokensTruncated returns the truncated shares of a delegation given
// a bond amount. It returns an error if the defi has no tokens.
func (v Defi) SharesFromTokensTruncated(amt sdk.Int) (sdk.Dec, error) {
	if v.Tokens.IsZero() {
		return sdk.ZeroDec(), ErrInsufficientShares
	}

	return v.GetDelegatorShares().MulInt(amt).QuoTruncate(v.GetTokens().ToDec()), nil
}

// get the bonded tokens which the defi holds
func (v Defi) BondedTokens() sdk.Int {
	if v.IsBonded() {
		return v.Tokens
	}

	return sdk.ZeroInt()
}

// UpdateStatus updates the location of the shares within a defi
// to reflect the new status
func (v Defi) UpdateStatus(newStatus BondStatus) Defi {
	v.Status = newStatus
	return v
}

// AddTokensFromDel adds tokens to a defi
func (v Defi) AddTokensFromDel(amount sdk.Int) (Defi, sdk.Dec) {
	// calculate the shares to issue
	var issuedShares sdk.Dec
	if v.DelegatorShares.IsZero() {
		// the first delegation to a defi sets the exchange rate to one
		issuedShares = amount.ToDec()
	} else {
		shares, err := v.SharesFromTokens(amount)
		if err != nil {
			panic(err)
		}

		issuedShares = shares
	}

	v.Tokens = v.Tokens.Add(amount)
	v.DelegatorShares = v.DelegatorShares.Add(issuedShares)

	return v, issuedShares
}

// RemoveTokens removes tokens from a defi
func (v Defi) RemoveTokens(tokens sdk.Int) Defi {
	if tokens.IsNegative() {
		panic(fmt.Sprintf("should not happen: trying to remove negative tokens %v", tokens))
	}

	if v.Tokens.LT(tokens) {
		panic(fmt.Sprintf("should not happen: only have %v tokens, trying to remove %v", v.Tokens, tokens))
	}

	v.Tokens = v.Tokens.Sub(tokens)

	return v
}

// RemoveDelShares removes delegator shares from a defi.
// NOTE: because token fractions are left in the defi,
//       the exchange rate of future shares of this defi can increase.
func (v Defi) RemoveDelShares(delShares sdk.Dec) (Defi, sdk.Int) {
	remainingShares := v.DelegatorShares.Sub(delShares)

	var issuedTokens sdk.Int
	if remainingShares.IsZero() {
		// last delegation share gets any trimmings
		issuedTokens = v.Tokens
		v.Tokens = sdk.ZeroInt()
	} else {
		// leave excess tokens in the defi
		// however fully use all the delegator shares
		issuedTokens = v.TokensFromShares(delShares).TruncateInt()
		v.Tokens = v.Tokens.Sub(issuedTokens)

		if v.Tokens.IsNegative() {
			panic("attempting to remove more tokens than available in defi")
		}
	}

	v.DelegatorShares = remainingShares

	return v, issuedTokens
}

// MinEqual defines a more minimum set of equality conditions when comparing two
// defis.
func (v *Defi) MinEqual(other *Defi) bool {
	return v.OperatorAddress == other.OperatorAddress &&
		v.Status == other.Status &&
		v.Tokens.Equal(other.Tokens) &&
		v.DelegatorShares.Equal(other.DelegatorShares) &&
		v.Description.Equal(other.Description) &&
		v.MinSelfDelegation.Equal(other.MinSelfDelegation)

}

// Equal checks if the receiver equals the parameter
func (v *Defi) Equal(v2 *Defi) bool {
	return v.MinEqual(v2) &&
		v.UnbondingHeight == v2.UnbondingHeight &&
		v.UnbondingTime.Equal(v2.UnbondingTime)
}

func (v Defi) GetMoniker() string    { return v.Description.Moniker }
func (v Defi) GetStatus() BondStatus { return v.Status }
func (v Defi) GetOperator() sdk.ValAddress {
	if v.OperatorAddress == "" {
		return nil
	}
	addr, err := sdk.ValAddressFromBech32(v.OperatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

func (v Defi) GetTokens() sdk.Int            { return v.Tokens }
func (v Defi) GetBondedTokens() sdk.Int      { return v.BondedTokens() }
func (v Defi) GetMinSelfDelegation() sdk.Int { return v.MinSelfDelegation }
func (v Defi) GetDelegatorShares() sdk.Dec   { return v.DelegatorShares }

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (v Defi) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return nil
}

// -----------------------------------------------------
// Rewards
// create a new DefiHistoricalRewards
func NewDefiHistoricalRewards(cumulativeRewardRatio sdk.DecCoins, referenceCount uint32) DefiHistoricalRewards {
	return DefiHistoricalRewards{
		CumulativeRewardRatio: cumulativeRewardRatio,
		ReferenceCount:        referenceCount,
	}
}

// create a new DefiCurrentRewards
func NewDefiCurrentRewards(rewards sdk.DecCoins, period uint64) DefiCurrentRewards {
	return DefiCurrentRewards{
		Rewards: rewards,
		Period:  period,
	}
}

// return the initial accumulated commission (zero)
func InitialDefiAccumulatedCommission() DefiAccumulatedCommission {
	return DefiAccumulatedCommission{}
}

