package store

import (
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"github.com/jorgepiresg/ChallangePismo/store/accounts"
	operationsType "github.com/jorgepiresg/ChallangePismo/store/operations_type"
	"github.com/jorgepiresg/ChallangePismo/store/transactions"
)

type Store struct {
	Accounts       accounts.IAccounts
	Transactions   transactions.ITransactions
	OperationsType operationsType.IOperationsType
}

type Options struct {
	DB    *sqlx.DB
	Log   *logrus.Logger
	Cache *redis.Client
}

func New(opts Options) Store {
	accountsOpts := accounts.Options{
		DB:    opts.DB,
		Log:   opts.Log,
		Cache: opts.Cache,
	}

	transactionsOpts := transactions.Options{
		DB:  opts.DB,
		Log: opts.Log,
	}

	operationsTypeOpts := operationsType.Options{
		DB:    opts.DB,
		Log:   opts.Log,
		Cache: opts.Cache,
	}

	return Store{
		Accounts:       accounts.New(accountsOpts),
		Transactions:   transactions.New(transactionsOpts),
		OperationsType: operationsType.New(operationsTypeOpts),
	}
}
