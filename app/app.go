package app

import (
	"github.com/jorgepiresg/ChallangePismo/app/accounts"
	"github.com/jorgepiresg/ChallangePismo/app/transactions"
	"github.com/jorgepiresg/ChallangePismo/store"
)

type App struct {
	Accounts     accounts.IAccounts
	Transactions transactions.ITransactions
}

type Options struct {
	Store store.Store
}

func New(opts Options) App {
	return App{
		Accounts:     accounts.New(accounts.Options{Store: opts.Store}),
		Transactions: transactions.New(transactions.Options{Store: opts.Store}),
	}
}
