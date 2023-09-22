package transactions

import (
	"context"
	"fmt"

	modelTransactions "github.com/jorgepiresg/ChallangePismo/model/transactions"
	"github.com/jorgepiresg/ChallangePismo/store"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/app/transactions_mock.go -package=mocksApp
type ITransactions interface {
	Make(ctx context.Context, data modelTransactions.MakeTransaction) error
}

type Options struct {
	Store store.Store
}

type transactions struct {
	store store.Store
}

func New(opts Options) ITransactions {
	return transactions{
		store: opts.Store,
	}
}

func (t transactions) Make(ctx context.Context, data modelTransactions.MakeTransaction) error {

	if err := data.ValidateAmount(); err != nil {
		return err
	}

	operationType, err := t.store.OperationsType.GetByID(ctx, data.OperationTypeID)
	if err != nil {
		return fmt.Errorf("operation type id not found")
	}

	if _, err := t.store.Accounts.GetByID(ctx, data.AccountID); err != nil {
		return fmt.Errorf("account id not found")
	}

	data.SetOperationInAmount(operationType.Operation)

	if _, err := t.store.Transactions.Create(ctx, data); err != nil {
		return fmt.Errorf("fail to make transaction")
	}

	return nil
}
