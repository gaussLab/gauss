package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/tmhash"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	emptyAddr string

	addr1 = sdk.AccAddress(tmhash.SumTruncated([]byte("addr1"))).String()
	addr2 = sdk.AccAddress(tmhash.SumTruncated([]byte("addr2"))).String()
)

// test ValidateBasic for MsgIssueToken
func TestMsgIssueAsset(t *testing.T) {
	addr := sdk.AccAddress(tmhash.SumTruncated([]byte("test"))).String()

	tests := []struct {
		testCase string
		*MsgIssueToken
		expectPass bool
	}{
		{"token unlocked", NewMsgIssueToken("Gauss Network", "stake", "ustake", 6, 1, 1, true, true, addr), true},
		{"token locked", NewMsgIssueToken("Gauss Network", "stake", "ustake", 6, 1, 1, true, false, addr), true},
		{"symbol empty", NewMsgIssueToken("Gauss Network", "", "", 6, 1, 1, true, true, addr), false},
		{"symbol error", NewMsgIssueToken("Gauss Network", "b&stake", "ub&stake", 6, 1, 1, true, true, addr), false},
		{"symbol first letter is num", NewMsgIssueToken("Gauss Network", "4stake", "u4stake", 6, 1, 1, true, true, addr), false},
		{"symbol too long", NewMsgIssueToken("Gauss Network", "stake123456789012345678901234567890123456789012345678901234567890", "ustake", 6, 1, 1, true, true, addr), false},
		{"unit too long", NewMsgIssueToken("Gauss Network", "stake", "ustake123456789012345678901234567890123456789012345678901234567890", 6, 1, 1, true, true, addr), false},
		{"symbol too short", NewMsgIssueToken("Gauss Network", "aa", "uaa", 6, 1, 1, true, true, addr), false},
		{"name empty", NewMsgIssueToken("", "stake", "ustake", 6, 1, 1, true, true, addr), false},
		{"name too long", NewMsgIssueToken("Gauss Network aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "stake", "ustake", 6, 1, 1, true, true, addr), false},
		{"initial supply is zero", NewMsgIssueToken("Gauss Network", "stake", "ustake", 6, 0, 1, true, true, addr), true},
		{"total supply is zero", NewMsgIssueToken("Gauss Network", "stake", "ustake", 6, 1, 0, true, true, addr), true},
		{"initial supply bigger than total supply", NewMsgIssueToken("Gauss Network", "stake", "ustake", 6, 2, 1, true, true, addr), false},
		{"decimals error", NewMsgIssueToken("Gauss Network", "stake", "ustake", 20, 1, 1, true, true, addr), false},
	}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, tc.MsgIssueToken.ValidateBasic(), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, tc.MsgIssueToken.ValidateBasic(), "test: %v", tc.testCase)
		}
	}
}

// test ValidateBasic for MsgEditToken
func TestMsgEditToken(t *testing.T) {
	owner := sdk.AccAddress(tmhash.SumTruncated([]byte("owner"))).String()
	mintable := False

	tests := []struct {
		testCase string
		*MsgEditToken
		expectPass bool
	}{
		{"basic good", NewMsgEditToken("ttk",  mintable, owner), true},
		{"symbol error", NewMsgEditToken("tt", mintable, ""), false},
		{"loss owner", NewMsgEditToken("ttk",  mintable, ""), false},
	}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, tc.MsgEditToken.ValidateBasic(), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, tc.MsgEditToken.ValidateBasic(), "test: %v", tc.testCase)
		}
	}
}

func TestMsgEditTokenRoute(t *testing.T) {
	symbol := "ttk"
	mintable := False

	// build a MsgEditToken
	msg := MsgEditToken{
		Name: "Test Token",
		Symbol:    symbol,
		Mintable:  mintable,
	}

	require.Equal(t, "token", msg.Route())
}

func TestMsgEditTokenGetSignBytes(t *testing.T) {
	mintable := False

	var msg = MsgEditToken{
		Name:     "Test Token",
		Symbol:    "ttk",
		Owner:     sdk.AccAddress(tmhash.SumTruncated([]byte("owner"))).String(),
		Mintable:  mintable,
	}

	res := msg.GetSignBytes()

	expected := `{"type":"gauss/token/MsgEditToken",
		"value":{"mintable":"false","name":"Test ERC20",
		"owner":"gauss1fsgzj6t7udv8zhf6zj32mkqhcjcpv52ygswxa5","symbol":"kkt"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgMintTokenValidateBasic(t *testing.T) {
	testData := []struct {
		testCase   string
		symbol     string
		owner      string
		to         string
		amount     uint64
		expectPass bool
	}{
		{"basic good", "kkt", addr1, addr2, 1000, true},
		{"symbol empty", "", addr1, addr2, 1000, false},
		{"symbol too short", "kk", addr1, addr2, 1000, false},
		{"owner empty", "kkt", emptyAddr, addr2, 1000, false},
		{"to empty", "kkt", addr1, emptyAddr, 1000, true},
		{"amount invalid", "btc", addr1, addr2, 0, false},
	}

	for _, td := range testData {
		msg := NewMsgMintToken(td.symbol, td.owner, td.to, td.amount)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.testCase)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.testCase)
		}
	}
}

func TestMsgBurnTokenValidateBasic(t *testing.T) {
	testData := []struct {
		testCase   string
		symbol     string
		sender     string
		amount     uint64
		expectPass bool
	}{
		{"basic good", "kkt", addr1, 1000, true},
		{"symbol empty", "", addr1, 1000, false},
		{"symbol too short", "kk", addr1, 1000, false},
		{"sender empty", "kkt", emptyAddr, 1000, false},
		{"amount invalid", "kkt", addr1, 0, false},
	}

	for _, td := range testData {
		msg := NewMsgBurnToken(td.symbol, td.sender, td.amount)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.testCase)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.testCase)
		}
	}
}

func TestMsgUnlockTokenValidateBasic(t *testing.T) {
	testData := []struct {
		testCase   string
		symbol     string
		owner      string
		expectPass bool
	}{
		{"basic good", "kkt", addr1, true},
		{"symbol empty", "", addr1, false},
		{"symbol too short", "kk", addr1, false},
		{"sender empty", "kkt", emptyAddr,false},
		{"amount invalid", "kkt", addr1, false},
	}

	for _, td := range testData {
		msg := NewMsgUnlockToken(td.symbol, td.owner)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.testCase)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.testCase)
		}
	}

}

func TestMsgTransferTokenOwnerValidation(t *testing.T) {
	testData := []struct {
		testCase   string
		oldOwner   string
		symbol     string
		newOwner   string
		expectPass bool
	}{
		{"basic good", addr1, "kkt", addr2, true},
		{"oldOwner empty", emptyAddr, "kkt", addr1, false},
		{"newOwner empty", addr1, "kkt", emptyAddr, false},
		{"symbol empty", addr1, "", addr2, false},
		{"symbol error", addr1, "kkt_min", addr2, false},
	}

	for _, td := range testData {
		msg := NewMsgTransferTokenOwner(td.symbol, td.oldOwner, td.newOwner)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.testCase)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.testCase)
		}
	}
}
