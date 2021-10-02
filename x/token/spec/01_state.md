<!--
order: 1
-->

# State

## Token

Definition of data structure of FungibleToken

- Token: `0x1 -> amino(Token)`

```go
type Token struct {
  Name          string
  Symbol        string
  Decimals      uint32
  InitialSupply uint64
  TotalSupply   uint64
  Mintable      bool
  Owner         string
}
```

## Params

Params is a module-wide configuration structure that stores system
parameters and defines overall functioning of the token module.

- Params: `Paramsspace("token") -> amino(params)`

```go
type Params struct {
  IssueTokenMinFee sdk.Coin
}
```

