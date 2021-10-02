package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the defi module
	ModuleName = "defi"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the defi module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the defi module
	RouterKey = ModuleName
)

var (
	FeePoolKey                        = []byte{0x11} // key for global defi state

	// Keys for store prefixes
	DefisKey                          = []byte{0x21} // prefix for each key to a defi

	DelegationKey                     = []byte{0x31} // key for a delegation
	UnbondingDelegationKey            = []byte{0x32} // key for an unbonding-delegation
	UnbondingDelegationByDefiIndexKey = []byte{0x33} // prefix for each key for an unbonding-delegation, by defi operator

	UnbondingQueueKey    = []byte{0x41} // prefix for the timestamps in unbonding queue
	DefiQueueKey         = []byte{0x43} // prefix for the timestamps in defi queue

	HistoricalInfoKey    = []byte{0x50} // prefix for the historical info

	DefiOutstandingRewardsPrefix    = []byte{0x61} // key for outstanding rewards
	DelegatorWithdrawAddrPrefix     = []byte{0x62} // key for delegator withdraw address
	DelegatorStartingInfoPrefix     = []byte{0x63} // key for delegator starting info
	DefiHistoricalRewardsPrefix     = []byte{0x64} // key for historical defis rewards / stake
	DefiCurrentRewardsPrefix        = []byte{0x65} // key for current defi rewards
	DefiAccumulatedCommissionPrefix = []byte{0x66} // key for accumulated defi commission
)

// gets the key for the defi with address
// VALUE: defi/Defi
func GetDefiKey(operatorAddr sdk.ValAddress) []byte {
	return append(DefisKey, operatorAddr.Bytes()...)
}

// GetDefiQueueKey returns the prefix key used for getting a set of unbonding
// defis whose unbonding completion occurs at the given time and height.
func GetDefiQueueKey(timestamp time.Time, height int64) []byte {
	heightBz := sdk.Uint64ToBigEndian(uint64(height))
	timeBz := sdk.FormatTimeBytes(timestamp)
	timeBzL := len(timeBz)
	prefixL := len(DefiQueueKey)

	bz := make([]byte, prefixL+8+timeBzL+8)

	// copy the prefix
	copy(bz[:prefixL], DefiQueueKey)

	// copy the encoded time bytes length
	copy(bz[prefixL:prefixL+8], sdk.Uint64ToBigEndian(uint64(timeBzL)))

	// copy the encoded time bytes
	copy(bz[prefixL+8:prefixL+8+timeBzL], timeBz)

	// copy the encoded height
	copy(bz[prefixL+8+timeBzL:], heightBz)

	return bz
}

// ParseDefiQueueKey returns the encoded time and height from a key created
// from GetDefiQueueKey.
func ParseDefiQueueKey(bz []byte) (time.Time, int64, error) {
	prefixL := len(DefiQueueKey)
	if prefix := bz[:prefixL]; !bytes.Equal(prefix, DefiQueueKey) {
		return time.Time{}, 0, fmt.Errorf("invalid prefix; expected: %X, got: %X", DefiQueueKey, prefix)
	}

	timeBzL := sdk.BigEndianToUint64(bz[prefixL : prefixL+8])
	ts, err := sdk.ParseTimeBytes(bz[prefixL+8 : prefixL+8+int(timeBzL)])
	if err != nil {
		return time.Time{}, 0, err
	}

	height := sdk.BigEndianToUint64(bz[prefixL+8+int(timeBzL):])

	return ts, int64(height), nil
}

// gets the key for delegator bond with defi
// VALUE: defi/Delegation
func GetDelegationKey(delAddr sdk.AccAddress, defiAddr sdk.ValAddress) []byte {
	return append(GetDelegationsKey(delAddr), defiAddr.Bytes()...)
}

// gets the prefix for a delegator for all defis
func GetDelegationsKey(delAddr sdk.AccAddress) []byte {
	return append(DelegationKey, delAddr.Bytes()...)
}

// gets the key for an unbonding delegation by delegator and defi addr
// VALUE: defi/UnbondingDelegation
func GetUBDKey(delAddr sdk.AccAddress, defiAddr sdk.ValAddress) []byte {
	return append(
		GetUBDsKey(delAddr.Bytes()),
		defiAddr.Bytes()...)
}

// gets the index-key for an unbonding delegation, stored by defi-index
// VALUE: none (key rearrangement used)
func GetUBDByValIndexKey(delAddr sdk.AccAddress, defiAddr sdk.ValAddress) []byte {
	return append(GetUBDsByDefiIndexKey(defiAddr), delAddr.Bytes()...)
}

// rearranges the DefiIndexKey to get the UBDKey
func GetUBDKeyFromDefiIndexKey(indexKey []byte) []byte {
	addrs := indexKey[1:] // remove prefix bytes
	if len(addrs) != 2*sdk.AddrLen {
		panic("unexpected key length")
	}

	defiAddr := addrs[:sdk.AddrLen]
	delAddr := addrs[sdk.AddrLen:]

	return GetUBDKey(delAddr, defiAddr)
}

// gets the prefix for all unbonding delegations from a delegator
func GetUBDsKey(delAddr sdk.AccAddress) []byte {
	return append(UnbondingDelegationKey, delAddr.Bytes()...)
}

// gets the prefix keyspace for the indexes of unbonding delegations for a defi
func GetUBDsByDefiIndexKey(defiAddr sdk.ValAddress) []byte {
	return append(UnbondingDelegationByDefiIndexKey, defiAddr.Bytes()...)
}

// gets the prefix for all unbonding delegations from a delegator
func GetUnbondingDelegationTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(UnbondingQueueKey, bz...)
}

// GetHistoricalInfoKey returns a key prefix for indexing HistoricalInfo objects.
func GetHistoricalInfoKey(height int64) []byte {
	return append(HistoricalInfoKey, []byte(strconv.FormatInt(height, 10))...)
}


// -----------------------------
// gets an address from a defi's outstanding rewards key
func GetDefiOutstandingRewardsAddress(key []byte) (defiAddr sdk.ValAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.ValAddress(addr)
}

// gets an address from a delegator's withdraw info key
func GetDelegatorWithdrawInfoAddress(key []byte) (delAddr sdk.AccAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.AccAddress(addr)
}

// gets the addresses from a delegator starting info key
func GetDelegatorStartingInfoAddresses(key []byte) (defiAddr sdk.ValAddress, delAddr sdk.AccAddress) {
	addr := key[1 : 1+sdk.AddrLen]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	defiAddr = sdk.ValAddress(addr)
	addr = key[1+sdk.AddrLen:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	delAddr = sdk.AccAddress(addr)
	return
}

// gets the address & period from a defi's historical rewards key
func GetDefiHistoricalRewardsAddressPeriod(key []byte) (defiAddr sdk.ValAddress, period uint64) {
	addr := key[1 : 1+sdk.AddrLen]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	defiAddr = sdk.ValAddress(addr)
	b := key[1+sdk.AddrLen:]
	if len(b) != 8 {
		panic("unexpected key length")
	}
	period = binary.LittleEndian.Uint64(b)
	return
}

// gets the address from a defi's current rewards key
func GetDefiCurrentRewardsAddress(key []byte) (defiAddr sdk.ValAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.ValAddress(addr)
}

// gets the address from a defi's accumulated commission key
func GetDefiAccumulatedCommissionAddress(key []byte) (defiAddr sdk.ValAddress) {
	addr := key[1:]
	if len(addr) != sdk.AddrLen {
		panic("unexpected key length")
	}
	return sdk.ValAddress(addr)
}


// gets the outstanding rewards key for a defi
func GetDefiOutstandingRewardsKey(valAddr sdk.ValAddress) []byte {
	return append(DefiOutstandingRewardsPrefix, valAddr.Bytes()...)
}

// gets the key for a delegator's withdraw addr
func GetDelegatorWithdrawAddrKey(delAddr sdk.AccAddress) []byte {
	return append(DelegatorWithdrawAddrPrefix, delAddr.Bytes()...)
}

// gets the key for a delegator's starting info
func GetDelegatorStartingInfoKey(v sdk.ValAddress, d sdk.AccAddress) []byte {
	return append(append(DelegatorStartingInfoPrefix, v.Bytes()...), d.Bytes()...)
}

// gets the prefix key for a defi's historical rewards
func GetDefiHistoricalRewardsPrefix(v sdk.ValAddress) []byte {
	return append(DefiHistoricalRewardsPrefix, v.Bytes()...)
}

// gets the key for a defi's historical rewards
func GetDefiHistoricalRewardsKey(v sdk.ValAddress, k uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, k)
	return append(append(DefiHistoricalRewardsPrefix, v.Bytes()...), b...)
}

// gets the key for a defi's current rewards
func GetDefiCurrentRewardsKey(v sdk.ValAddress) []byte {
	return append(DefiCurrentRewardsPrefix, v.Bytes()...)
}

// gets the key for a validator's current commission
func GetDefiAccumulatedCommissionKey(v sdk.ValAddress) []byte {
	return append(DefiAccumulatedCommissionPrefix, v.Bytes()...)
}

