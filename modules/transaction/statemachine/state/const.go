package state

const (
	selectWalletForUpdate = `SELECT wallet_id, amount FROM wallet_assets WHERE wallet_id = ? AND currency = ? FOR UPDATE`
	updateAmountOfWallet  = `UPDATE wallet_assets SET amount = ? WHERE wallet_id = ? AND currency = ?`
	selectCard            = `SELECT wallet_id FROM cards WHERE id = ?`
)
