package transactions

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	mocksStore "github.com/jorgepiresg/ChallangePismo/mocks/store"
	modelAccounts "github.com/jorgepiresg/ChallangePismo/model/accounts"
	modelOperaTionsType "github.com/jorgepiresg/ChallangePismo/model/operations_type"
	modelTransactions "github.com/jorgepiresg/ChallangePismo/model/transactions"
	"github.com/jorgepiresg/ChallangePismo/store"
)

func TestMake(t *testing.T) {

	type fields struct {
		transactions   *mocksStore.MockITransactions
		accounts       *mocksStore.MockIAccounts
		operationsType *mocksStore.MockIOperationsType
	}

	tests := map[string]struct {
		input   modelTransactions.MakeTransaction
		err     error
		prepare func(f *fields)
	}{
		"should be able to make a new transaction": {
			input: modelTransactions.MakeTransaction{
				AccountID:       "id",
				OperationTypeID: 1,
				Amount:          10.00,
			},
			prepare: func(f *fields) {
				f.operationsType.EXPECT().GetByID(gomock.Any(), 1).Times(1).Return(modelOperaTionsType.OperationType{
					OperationTypeID: 1,
					Description:     "COMPRA A VISTA",
					Operation:       -1,
				}, nil)

				f.accounts.EXPECT().GetByID(gomock.Any(), "id").Times(1).Return(modelAccounts.Account{ID: "id"}, nil)

				f.transactions.EXPECT().Create(gomock.Any(), modelTransactions.MakeTransaction{
					AccountID:       "id",
					Amount:          -10.00,
					OperationTypeID: 1,
				}).Times(1).Return(modelTransactions.Transaction{}, nil)
			},
		},
		"should not be able to make a new transaction with error amount is invalid": {
			input: modelTransactions.MakeTransaction{
				Amount: 10.001,
			},
			prepare: func(f *fields) {},
			err:     fmt.Errorf("amount %v is invalid, use 2 decimals", 10.001),
		},
		"should not be able to make a new transaction with error amount 0 is invalid": {
			input: modelTransactions.MakeTransaction{
				Amount: 0,
			},
			prepare: func(f *fields) {},
			err:     fmt.Errorf("amount 0 is invalid"),
		},
		"should not be able to make a new transaction with error operation type id not found": {
			input: modelTransactions.MakeTransaction{
				Amount: 10,
			},
			prepare: func(f *fields) {
				f.operationsType.EXPECT().GetByID(gomock.Any(), 0).Times(1).Return(modelOperaTionsType.OperationType{}, fmt.Errorf("any"))
			},
			err: fmt.Errorf("operation type id not found"),
		},
		"should not be able to make a new transaction with error account id not found": {
			input: modelTransactions.MakeTransaction{
				AccountID:       "invalid_id",
				OperationTypeID: 1,
				Amount:          10.00,
			},
			prepare: func(f *fields) {
				f.operationsType.EXPECT().GetByID(gomock.Any(), 1).Times(1).Return(modelOperaTionsType.OperationType{
					OperationTypeID: 1,
					Description:     "COMPRA A VISTA",
					Operation:       -1,
				}, nil)

				f.accounts.EXPECT().GetByID(gomock.Any(), "invalid_id").Times(1).Return(modelAccounts.Account{}, fmt.Errorf("any"))
			},
			err: fmt.Errorf("account id not found"),
		},
		"should not be able to make a new transaction with error fail to make transaction": {
			input: modelTransactions.MakeTransaction{
				AccountID:       "id",
				OperationTypeID: 1,
				Amount:          10.00,
			},
			prepare: func(f *fields) {
				f.operationsType.EXPECT().GetByID(gomock.Any(), 1).Times(1).Return(modelOperaTionsType.OperationType{
					OperationTypeID: 1,
					Description:     "COMPRA A VISTA",
					Operation:       -1,
				}, nil)

				f.accounts.EXPECT().GetByID(gomock.Any(), "id").Times(1).Return(modelAccounts.Account{ID: "id"}, nil)

				f.transactions.EXPECT().Create(gomock.Any(), modelTransactions.MakeTransaction{
					AccountID:       "id",
					Amount:          -10.00,
					OperationTypeID: 1,
				}).Times(1).Return(modelTransactions.Transaction{}, fmt.Errorf("any"))
			},
			err: fmt.Errorf("fail to make transaction"),
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			accountsMock := mocksStore.NewMockIAccounts(ctrl)
			transactionsMock := mocksStore.NewMockITransactions(ctrl)
			operationsTypeMock := mocksStore.NewMockIOperationsType(ctrl)

			tt.prepare(&fields{
				accounts:       accountsMock,
				transactions:   transactionsMock,
				operationsType: operationsTypeMock,
			})

			a := New(Options{
				Store: store.Store{
					Accounts:       accountsMock,
					Transactions:   transactionsMock,
					OperationsType: operationsTypeMock,
				},
			})

			err := a.Make(context.Background(), tt.input)
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf(`Expected err: "%s" got "%s"`, tt.err, err)
			}
		})
	}
}
