package types

const (
	DepositType  Type = "deposit"
	WithdrawType Type = "withdraw"
	TransferType Type = "transfer"

	ValidatingState               State = "validate"
	WithDrawFromSenderState       State = "withdraw_from_sender"
	WithDrawFromBankState         State = "withdraw_from_bank"
	DepositToReceiverState        State = "deposit_to_receiver"
	DepositToBankState            State = "deposit_to_bank"
	VerifyBankingTransactionState State = "verify_banking_transaction"
	RefundToSenderState           State = "refund_to_sender"
	CompletedState                State = "complete"
	RejectedState                 State = "reject"

	CompleteStatus Status = "complete"
	SuccessStatus  Status = "success"
	FailStatus     Status = "failed"
) 
