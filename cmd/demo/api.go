package main

import (
	"context"
	"encoding/json"
	"ewallet/modules/asset"
	"ewallet/modules/wallet"
	"fmt"
	"io"
	"log"
	"net/http"
)
import "github.com/gorilla/mux"

type (
	apiHandler struct {
		service wallet.Service
	}
)

func startServer(service wallet.Service)  {
	handler := apiHandler{
		service: service,
	}
	r := mux.NewRouter()

	r.HandleFunc("/balance/{wallet_id}", handler.Balance).Methods(http.MethodGet)
	r.HandleFunc("/transfer", handler.Transfer).Methods(http.MethodPost)
	r.HandleFunc("/deposit", handler.Deposit).Methods(http.MethodPost)
	r.HandleFunc("/withdraw", handler.Withdraw).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func (receiver *apiHandler) Balance(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	assets, err := receiver.service.BalanceOf(context.TODO(), wallet.Wallet{
		ID: vars["wallet_id"],
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, _ := json.Marshal(assets)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(result))
}

func (receiver *apiHandler) Transfer(w http.ResponseWriter, r *http.Request)  {
	var dto TransferDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fromWallet := wallet.Wallet{
		ID: dto.From,
	}
	toWallet := wallet.Wallet{
		ID: dto.To,
	}
	money := asset.Asset{
		Amount: dto.Amount,
		Currency: dto.Currency,
	}
	transID, err := receiver.service.Send(context.TODO(), fromWallet, toWallet, money)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, fmt.Sprintf(`{"transaction_id":%d}`, transID))
}

func (receiver *apiHandler) Deposit(w http.ResponseWriter, r *http.Request)  {
	return
}

func (receiver *apiHandler) Withdraw(w http.ResponseWriter, r *http.Request)  {
	return
}


