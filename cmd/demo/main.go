package main

import (
	"context"
	"ewallet/modules/share"
	"ewallet/modules/transaction"
	"ewallet/modules/wallet"
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
	walletDB, err := sqlx.Connect("mysql", conf.WalletDatasource)
	if err != nil {
		panic(err)
	}
	
	transIDChannel := make(chan int64, 100)
	publisher := NewPublisher(transIDChannel)
	consumer := NewConsumer(transIDChannel, buildTransactionHandler(transDB, walletDB))

	walletRepo := wallet.NewRepository(walletDB)
	transactionService := buildTransactionService(conf.TransactionDatasource, walletRepo, publisher)
	walletService := wallet.NewService(walletRepo, transactionService)
	
	consumer.Run(context.Background())

	startServer(walletService)
}

func buildTransactionService(datasource string, walletRepo transaction.WalletRepository, publisher transaction.Publisher) transaction.Service {
	transDB, err := sqlx.Connect("mysql", datasource)
	if err != nil {
		panic(err)
	}
	return transaction.NewService(transaction.NewRepository(transDB, publisher), walletRepo)
}

func buildTransactionHandler(transactionDB ,assetDB share.DB) Handler {
	return transaction.NewStateMachine(transactionDB, assetDB, transaction.NewStateHandlerGetter(assetDB))
}