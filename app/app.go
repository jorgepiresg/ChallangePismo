package app

import (
	"log"

	"github.com/jorgepiresg/ChallangePismo/app/accounts"
	"github.com/jorgepiresg/ChallangePismo/app/transactions"
	"github.com/jorgepiresg/ChallangePismo/store"
	"github.com/sirupsen/logrus"
)

type App struct {
	Accounts     accounts.IAccounts
	Transactions transactions.ITransactions
}

type Options struct {
	Store store.Store
	Log   *logrus.Logger
}

func New(opts Options) App {
	app := App{
		Accounts:     accounts.New(accounts.Options{Store: opts.Store, Log: opts.Log}),
		Transactions: transactions.New(transactions.Options{Store: opts.Store, Log: opts.Log}),
	}

	log.Println("APP Created")
	return app
}
