package main

type (
	TransferDTO struct {
		From string `json:"from"`
		To string `json:"to"`
		Amount string `json:"amount"`
		Currency string `json:"currency"`
	}
)
