package transactions

import (
	"context"

	"github.com/jmoiron/sqlx"
	modelTransactions "github.com/jorgepiresg/ChallangePismo/model/transactions"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/store/transactions_mock.go -package=mocksStore
type ITransactions interface {
	Create(ctx context.Context, create modelTransactions.MakeTransaction) (modelTransactions.Transaction, error)
}

type Options struct {
	DB  *sqlx.DB
	Log *logrus.Logger
}

type transactions struct {
	db  *sqlx.DB
	log *logrus.Logger
}

func New(opts Options) ITransactions {
	return transactions{
		db:  opts.DB,
		log: opts.Log,
	}
}

func (t transactions) Create(ctx context.Context, create modelTransactions.MakeTransaction) (modelTransactions.Transaction, error) {

	var transaction modelTransactions.Transaction

	rows, err := t.db.NamedQueryContext(ctx, `INSERT INTO transactions (account_id, operation_type_id, amount) VALUES (:account_id, :operation_type_id, :amount) RETURNING *`, create)
	if err != nil {
		t.log.WithField("create", create).Error(err)
		return transaction, err
	}

	for rows.Next() {
		err = rows.StructScan(&transaction)
		if err != nil {
			t.log.WithField("create", create).Error(err)
			return transaction, err
		}
	}

	return transaction, nil
}
