package main

import (
	"context"
	service2 "ewallet/modules/bank/service"
	"ewallet/modules/share/dbclient"
	infra2 "ewallet/modules/transaction/infra"
	service3 "ewallet/modules/transaction/service"
	"ewallet/modules/transaction/statemachine/controller"
	"ewallet/modules/transaction/statemachine/state"
	"ewallet/modules/wallet/infra"
	"ewallet/modules/wallet/service"
	"fmt"
	"github.com/caarlos0/env"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	TransactionDatasource         string `env:"TRANSACTION_DATASOURCE" envDefault:"test:test@(localhost:3306)/transaction"`
	WalletDatasource         string `env:"WALLET_DATASOURCE" envDefault:"test:test@(localhost:3307)/wallet"`
}



func main()  {
	conf := Config{}
	if err := env.Parse(&conf); err != nil {
		panic(err.Error())
	}
	transDB, err := sqlx.Connect("mysql", conf.TransactionDatasource)
	if err != nil {
		panic(err)
	}
	defer transDB.Close()
	walletDB, err := sqlx.Connect("mysql", conf.WalletDatasource)
	if err != nil {
		panic(err)
	}
	defer transDB.Close()

	transferChannel := make(chan int64, 100)
	depositChannel := make(chan int64, 100)
	withdrawChannel := make(chan int64, 100)
	defer close(transferChannel)
	defer close(depositChannel)
	defer close(withdrawChannel)
	publisher := NewPublisher(transferChannel, depositChannel, withdrawChannel)
	transferConsumer := NewConsumer(transferChannel, buildTransferHandler(transDB, walletDB))
	depositConsumer := NewConsumer(depositChannel, buildDepositHandler(transDB, walletDB))
	withdrawConsumer := NewConsumer(withdrawChannel, buildWithdrawHandler(transDB, walletDB))

	walletRepo := infra.NewRepository(walletDB)
	transactionService := buildTransactionService(conf.TransactionDatasource, publisher)
	walletService := service.NewService(walletRepo, transactionService)

	transferConsumer.Run(context.Background())
	depositConsumer.Run(context.Background())
	withdrawConsumer.Run(context.Background())

	fmt.Println("running at localhost:8080")
	startServer(walletService)
}

func buildTransactionService(datasource string, publisher service3.Publisher) service3.Service {
	transDB, err := sqlx.Connect("mysql", datasource)
	if err != nil {
		panic(err)
	}
	return service3.NewService(
		infra2.NewRepository(transDB),
		publisher)
}

func buildTransferHandler(transactionDB ,assetDB dbclient.DB) Handler {
	bankService := service2.NewBank()
	return controller.NewStateMachineHandler(
		transactionDB,
		state.NewStateHandlerGetter(assetDB, bankService,  bankService,  bankService),
		controller.NewTransferStateController())
}

func buildDepositHandler(transactionDB ,assetDB dbclient.DB) Handler {
	bankService := service2.NewBank()
	return controller.NewStateMachineHandler(
		transactionDB,
		state.NewStateHandlerGetter(assetDB, bankService,  bankService,  bankService),
		controller.NewDepositStateController())
}

func buildWithdrawHandler(transactionDB ,assetDB dbclient.DB) Handler {
	bankService := service2.NewBank()
	return controller.NewStateMachineHandler(
		transactionDB,
		state.NewStateHandlerGetter(assetDB, bankService,  bankService,  bankService),
		controller.NewWithdrawStateController())
}