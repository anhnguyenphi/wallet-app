package asset

type (
	Asset struct {
		WalletID string `json:"wallet_id"`
		Currency string `json:"currency"`
		Amount string `json:"amount"`
	}
)