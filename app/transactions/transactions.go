package transactions

import (
	"context"
	"fmt"

	modelTransactions "github.com/jorgepiresg/ChallangePismo/model/transactions"
	"github.com/jorgepiresg/ChallangePismo/store"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/app/transactions_mock.go -package=mocksApp
type ITransactions interface {
	Make(ctx context.Context, data modelTransactions.MakeTransaction) error
}

type Options struct {
	Store store.Store
	Log   *logrus.Logger
}

type transactions struct {
	store store.Store
	log   *logrus.Logger
}

func New(opts Options) ITransactions {
	return transactions{
		store: opts.Store,
		log:   opts.Log,
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

	res, err := t.store.Transactions.Create(ctx, data)
	if err != nil {
		return fmt.Errorf("fail to make transaction")
	}

	go t.discharge(context.Background(), res)

	return nil
}

func (t transactions) discharge(ctx context.Context, data modelTransactions.Transaction) {

	if data.Amount <= 0 {
		return
	}

	transactions, err := t.store.Transactions.GetToDischargeByAccountID(ctx, data.AccountID)
	if err != nil {
		t.log.Error(err)
		return
	}

	if len(transactions) == 0 {
		return
	}

	currentBalance := data.Amount

	for _, transaction := range transactions {

		if currentBalance == 0 {
			break
		}

		currentBalance += transaction.Balance

		if currentBalance >= 0 {
			transaction.Balance = 0
		}

		if currentBalance < 0 {
			transaction.Balance = currentBalance
			currentBalance = 0
		}

		err := t.store.Transactions.UpdateBalance(ctx, transaction)
		if err != nil {
			t.log.Error(err)
			return
		}
	}

	data.Balance = currentBalance

	if err := t.store.Transactions.UpdateBalance(ctx, data); err != nil {
		t.log.Error(err)
	}
}
