package transactions

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	modelTransactions "github.com/jorgepiresg/ChallangePismo/model/transactions"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/store/transactions_mock.go -package=mocksStore
type ITransactions interface {
	Create(ctx context.Context, create modelTransactions.MakeTransaction) (modelTransactions.Transaction, error)
}

type transactions struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) ITransactions {
	return transactions{
		db: db,
	}
}

func (t transactions) Create(ctx context.Context, create modelTransactions.MakeTransaction) (modelTransactions.Transaction, error) {

	var transaction modelTransactions.Transaction

	rows, err := t.db.NamedQueryContext(ctx, `INSERT INTO transactions (account_id, operation_type_id, amount) VALUES (:account_id, :operation_type_id, :amount) RETURNING *`, create)
	if err != nil {
		fmt.Println(err.Error())
		return transaction, err
	}

	for rows.Next() {
		err = rows.StructScan(&transaction)
		if err != nil {
			return transaction, err
		}
	}

	return transaction, nil
}
