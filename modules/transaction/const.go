package transaction


const (
	ValidatingState         State = "validating"
	WithDrawFromSenderState State = "withdrawing_from_sender"
	DepositToReceiverState  State = "depositing_to_receiver"
	ReversingToSenderState  State = "reversing_to_sender"
	CompletedState          State = "completed"
	RejectedState           State = "rejected"

	selectForUpdateTransactionSql     = `SELECT * from transactions where id = ? FOR UPDATE`
	selectWalletForUpdate             = `SELECT * FROM assets WHERE wallet_id = ? AND currency = ? FOR UPDATE`
	updateAmountOfWallet              = `UPDATE assets SET amount = ? WHERE wallet_id = ? AND currency = ?`
	updateStateOfTransactionForUpdate = `UPDATE transactions SET state = ? WHERE id = ?`
	selectWallet = `SELECT * FROM assets WHERE wallet_id = ? AND currency = ?`
)
