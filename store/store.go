package store

import (
	"github.com/jmoiron/sqlx"

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
	DB *sqlx.DB
}

func New(opts Options) Store {
	return Store{
		Accounts:       accounts.New(opts.DB),
		Transactions:   transactions.New(opts.DB),
		OperationsType: operationsType.New(opts.DB),
	}
}
