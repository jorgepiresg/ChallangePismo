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
	GetToDischargeByAccountID(ctx context.Context, accountID string) ([]modelTransactions.Transaction, error)
	UpdateBalance(ctx context.Context, transaction modelTransactions.Transaction) error
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

	rows, err := t.db.NamedQueryContext(ctx, `INSERT INTO transactions (account_id, operation_type_id, amount, balance) VALUES (:account_id, :operation_type_id, :amount, :amount) RETURNING *`, create)
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

func (t transactions) GetToDischargeByAccountID(ctx context.Context, accountID string) ([]modelTransactions.Transaction, error) {

	var transactions []modelTransactions.Transaction
	err := t.db.SelectContext(ctx, &transactions, `SELECT transaction_id ,account_id, operation_type_id, amount, balance, event_date FROM transactions where 
	account_id = $1 AND
	balance < 0 
	ORDER BY event_date asc;
	`, accountID)

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (t transactions) UpdateBalance(ctx context.Context, transaction modelTransactions.Transaction) error {

	_, err := t.db.ExecContext(ctx, `UPDATE transactions SET balance = $1 where transaction_id = $2`, transaction.Balance, transaction.TransactionID)

	return err
}
