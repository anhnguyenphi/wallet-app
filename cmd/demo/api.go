package main

import (
	"context"
	"encoding/json"
	"ewallet/modules/asset"
	service2 "ewallet/modules/wallet/service"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)
import "github.com/gorilla/mux"

type (
	apiHandler struct {
		service service2.Service
	}
)

func startServer(service service2.Service)  {
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
	walletID, err := strconv.ParseInt(vars["wallet_id"], 10, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	assets, err := receiver.service.BalanceOf(context.TODO(), walletID)
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
	money := asset.Asset{
		Amount: dto.Amount,
		Currency: dto.Currency,
	}
	transID, err := receiver.service.Transfer(context.TODO(), dto.FromWalletID, dto.ToWalletID, money)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, fmt.Sprintf(`{"transaction_id":%d}`, transID))
}

func (receiver *apiHandler) Deposit(w http.ResponseWriter, r *http.Request)  {
	var dto DepositDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	money := asset.Asset{
		Amount: dto.Amount,
		Currency: dto.Currency,
	}
	transID, err := receiver.service.Deposit(context.TODO(), dto.FromCardID, dto.ToWalletID, money)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, fmt.Sprintf(`{"transaction_id":%d}`, transID))
}

func (receiver *apiHandler) Withdraw(w http.ResponseWriter, r *http.Request)  {
	var dto WithdrawDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	money := asset.Asset{
		Amount: dto.Amount,
		Currency: dto.Currency,
	}
	transID, err := receiver.service.Withdraw(context.TODO(), dto.FromWalletID, dto.ToCardID, money)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, fmt.Sprintf(`{"transaction_id":%d}`, transID))
}


