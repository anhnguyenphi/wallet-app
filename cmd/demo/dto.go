package main

type (
	TransferDTO struct {
		FromWalletID int64 `json:"from_wallet_id"`
		ToWalletID int64 `json:"to_wallet_id"`
		Amount string `json:"amount"`
		Currency string `json:"currency"`
	}

	DepositDTO struct {
		FromCardID int64 `json:"from_card_id"`
		ToWalletID int64 `json:"to_wallet_id"`
		Amount string `json:"amount"`
		Currency string `json:"currency"`
	}

	WithdrawDTO struct {
		FromWalletID int64 `json:"from_wallet_id"`
		ToCardID int64 `json:"to_card_id"`
		Amount string `json:"amount"`
		Currency string `json:"currency"`
	}
)
