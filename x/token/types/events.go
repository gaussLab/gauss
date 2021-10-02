package types

const (
	AttributeValueCategory = ModuleName

	EventTypeIssueToken         = "issue_token"
	EventTypeEditToken          = "edit_token"
	EventTypeMintToken          = "mint_token"
	EventTypeBurnToken          = "burn_token"
	EventTypeUnlockToken        = "unlock_token"
	EventTypeTransferTokenOwner = "transfer_token_owner"

	AttributeKeyCreator   = "creator"
	AttributeKeySymbol    = "symbol"
	AttributeKeyAmount    = "amount"
	AttributeKeyOwner     = "owner"
	AttributeKeyNewOwner  = "new_owner"
	AttributeKeyRecipient = "recipient"
)
