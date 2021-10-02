<!--
order: 2
-->

# Messages

In this section we describe the processing of the token messages and the
corresponding updates to the state.

## MsgIssueToken

A token is created using the `MsgIssueToken` message.

```go
type MsgIssueToken struct {
  Name          string
  Symbol        string
  Decimals      uint32
  InitialSupply uint64
  TotalSupply   uint64
  Mintable      bool
  Owner         string
}
```

This message is expected to fail if:

- the `Name` of the token is faulty, namely:
  - is not begin with `[a-zA-Z]`
  - contains characters other than letters and numbers
  - character length exceeds 32 bits
- the `Symbol` of the token is faulty, namely:
  - is not begin with `[a-zA-Z]`
  - contains characters other than letters and numbers
  - character length is greater than 8 bits or less than 3 bits
  - this symbol is already registered
- the `Decimals` > 18 or `Decimals` < 0
- the `TotalSupply` > `max` or `TotalSupply` < `InitialSupply`

This message creates and stores the `Token` object at appropriate
indexes.

## MsgEditToken

The `Name`, `TotalSupply`, `Mintable` of a token can be updated using the
`MsgEditToken`.

```go
type MsgEditToken struct {
  Symbol    string
  Name      string
  TotalSupply uint64
  Mintable  Bool
  Owner     string
}
```

This message is expected to fail if:

- the `Symbol` is not existed
- the `Name` of the token is faulty, namely:
  - is not begin with `[a-zA-Z]`
  - contains characters other than letters and numbers
  - character length exceeds 32 bits
- the `TotalSupply` > `max` or `TotalSupply` < `InitialSupply`
- the `Owner` is not the token owner

This message stores the updated `Token` object.

## MsgMintToken

The owner of the token can mint some tokens to the specified account

```go
type MsgMintToken struct {
  Symbol string
  Owner  string
  To     string
  Amount uint64
```

This message is expected to fail if:

- the `Symbol` is not existed
- the `Mintable` of the token is false
- the `Owner` is not the token owner
- the `Amount` `Coin` has exceeded the number of additional
  issuances（**TotalSupply - Issued**）

## MsgBurnToken

The owner of the token can mint some tokens to the specified account

```go
type MsgBurnToken struct {
  Symbol string
  Sender string
  Amount uint64
```

This message is expected to fail if:

- the `Symbol` is not existed
- the `Amount` don't have enough tokens

## MsgTransferTokenOwner

The ownership of the `token` can be transferred to others

```go
type MsgTransferTokenOwner struct {
  Symbol   string
  OldOwner string
  NewOwner string
}
```

This message is expected to fail if:

- the `Symbol` is not existed
- the `Owner` is not the token owner

