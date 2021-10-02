package types

// defi module event types
const (
	EventTypeCompleteUnbonding    = "complete_unbonding"
	EventTypeCreateDefi           = "create_defi"
	EventTypeEditDefi             = "edit_defi"
	EventTypeDelegate             = "delegate"
	EventTypeUnbond               = "unbond"
	EventTypeSetWithdrawAddress   = "set_withdraw_address"
	EventTypeRewards              = "rewards"
	EventTypeCommission           = "commission"
	EventTypeWithdrawRewards      = "withdraw_rewards"
	EventTypeWithdrawCommission   = "withdraw_commission"
	EventTypeWithdrawDelegatorRewards  = "withdraw_delegator_rewards"
	EventTypeMint                 = "mint"

	AttributeKeyDefi              = "defi"
	AttributeKeyWithdrawAddress   = "withdraw_address"
	AttributeKeyMinSelfDelegation = "min_self_delegation"
	AttributeKeyDelegator         = "delegator"
	AttributeKeyCompletionTime    = "completion_time"
	AttributeKeyRecipient         = "recipient"

	AttributeValueCategory        = ModuleName
)
