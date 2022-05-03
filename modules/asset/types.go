package asset

type (
	Asset struct {
		WalletID int64 `json:"wallet_id"`
		Currency string `json:"currency"`
		Amount string `json:"amount"`
	}
)