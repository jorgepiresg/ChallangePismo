package transactions

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	mocksStore "github.com/jorgepiresg/ChallangePismo/mocks/store"
	modelAccounts "github.com/jorgepiresg/ChallangePismo/model/accounts"
	modelOperaTionsType "github.com/jorgepiresg/ChallangePismo/model/operations_type"
	modelTransactions "github.com/jorgepiresg/ChallangePismo/model/transactions"
	"github.com/jorgepiresg/ChallangePismo/store"
	"github.com/sirupsen/logrus"
)

func TestMake(t *testing.T) {

	type fields struct {
		transactions   *mocksStore.MockITransactions
		accounts       *mocksStore.MockIAccounts
		operationsType *mocksStore.MockIOperationsType
		wg             *sync.WaitGroup
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
				Amount:          10.50,
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
					Amount:          -10.50,
					OperationTypeID: 1,
				}).Times(1).Return(modelTransactions.Transaction{}, nil)
			},
		},
		"should not be able to make a new transaction with error amount negative invalid": {
			input: modelTransactions.MakeTransaction{
				Amount: -10,
			},
			prepare: func(f *fields) {},
			err:     fmt.Errorf("amount invalid"),
		},
		"should not be able to make a new transaction with error amount 0 is invalid": {
			input: modelTransactions.MakeTransaction{
				Amount: 0,
			},
			prepare: func(f *fields) {},
			err:     fmt.Errorf("amount invalid"),
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
		"should be able to make a new transaction with dischard": {
			input: modelTransactions.MakeTransaction{
				AccountID:       "id",
				OperationTypeID: 4,
				Amount:          60.00,
			},
			prepare: func(f *fields) {

				f.operationsType.EXPECT().GetByID(gomock.Any(), 4).Times(1).Return(modelOperaTionsType.OperationType{
					OperationTypeID: 4,
					Description:     "PAGAMENTO",
					Operation:       1,
				}, nil)

				f.accounts.EXPECT().GetByID(gomock.Any(), "id").Times(1).Return(modelAccounts.Account{ID: "id"}, nil)

				f.transactions.EXPECT().Create(gomock.Any(), modelTransactions.MakeTransaction{
					AccountID:       "id",
					Amount:          60.00,
					OperationTypeID: 4,
				}).Times(1).Return(modelTransactions.Transaction{
					TransactionID:   "transaction_id",
					AccountID:       "id",
					Amount:          60.00,
					OperationTypeID: 4,
					Balance:         60,
				}, nil)

				f.wg.Add(4)

				f.transactions.EXPECT().GetToDischargeByAccountID(gomock.Any(), "id").Times(1).Return([]modelTransactions.Transaction{
					{
						TransactionID:   "1",
						AccountID:       "id",
						OperationTypeID: 1,
						Amount:          -50,
						Balance:         -50,
					},
					{
						TransactionID:   "2",
						AccountID:       "id",
						OperationTypeID: 1,
						Amount:          -23.50,
						Balance:         -23.50,
					},
				}, nil).Do(func(arg0, arg1 interface{}) {
					f.wg.Done()
				})

				f.transactions.EXPECT().UpdateBalance(gomock.Any(), modelTransactions.Transaction{
					TransactionID:   "1",
					AccountID:       "id",
					OperationTypeID: 1,
					Amount:          -50,
					Balance:         0,
				}).Times(1).Return(nil).Do(func(arg0, arg1 interface{}) {
					f.wg.Done()
				})

				f.transactions.EXPECT().UpdateBalance(gomock.Any(), modelTransactions.Transaction{
					TransactionID:   "2",
					AccountID:       "id",
					OperationTypeID: 1,
					Amount:          -23.50,
					Balance:         -13.50,
				}).Times(1).Return(nil).Do(func(arg0, arg1 interface{}) {
					f.wg.Done()
				})

				f.transactions.EXPECT().UpdateBalance(gomock.Any(), modelTransactions.Transaction{
					TransactionID:   "transaction_id",
					AccountID:       "id",
					OperationTypeID: 4,
					Amount:          60,
					Balance:         0,
				}).Times(1).Do(func(arg0, arg1 interface{}) {
					f.wg.Done()
				})

			},
		},

		"should be able to make a new transaction with dischard where current balance is zero": {
			input: modelTransactions.MakeTransaction{
				AccountID:       "id",
				OperationTypeID: 4,
				Amount:          60.00,
			},
			prepare: func(f *fields) {

				f.operationsType.EXPECT().GetByID(gomock.Any(), 4).Times(1).Return(modelOperaTionsType.OperationType{
					OperationTypeID: 4,
					Description:     "PAGAMENTO",
					Operation:       1,
				}, nil)

				f.accounts.EXPECT().GetByID(gomock.Any(), "id").Times(1).Return(modelAccounts.Account{ID: "id"}, nil)

				f.transactions.EXPECT().Create(gomock.Any(), modelTransactions.MakeTransaction{
					AccountID:       "id",
					Amount:          60.00,
					OperationTypeID: 4,
				}).Times(1).Return(modelTransactions.Transaction{
					TransactionID:   "transaction_id",
					AccountID:       "id",
					Amount:          60.00,
					OperationTypeID: 4,
					Balance:         60,
				}, nil)

				f.wg.Add(3)

				f.transactions.EXPECT().GetToDischargeByAccountID(gomock.Any(), "id").Times(1).Return([]modelTransactions.Transaction{
					{
						TransactionID:   "1",
						AccountID:       "id",
						OperationTypeID: 1,
						Amount:          -60,
						Balance:         -60,
					},
					{
						TransactionID:   "2",
						AccountID:       "id",
						OperationTypeID: 1,
						Amount:          -23.50,
						Balance:         -23.50,
					},
				}, nil).Do(func(arg0, arg1 interface{}) {
					f.wg.Done()
				})

				f.transactions.EXPECT().UpdateBalance(gomock.Any(), modelTransactions.Transaction{
					TransactionID:   "1",
					AccountID:       "id",
					OperationTypeID: 1,
					Amount:          -60,
					Balance:         0,
				}).Times(1).Return(nil).Do(func(arg0, arg1 interface{}) {
					f.wg.Done()
				})

				f.transactions.EXPECT().UpdateBalance(gomock.Any(), modelTransactions.Transaction{
					TransactionID:   "transaction_id",
					AccountID:       "id",
					OperationTypeID: 4,
					Amount:          60,
					Balance:         0,
				}).Times(1).Do(func(arg0, arg1 interface{}) {
					f.wg.Done()
				})

			},
		},

		"should be able to make a new transaction with error in dischard": {
			input: modelTransactions.MakeTransaction{
				AccountID:       "id",
				OperationTypeID: 4,
				Amount:          60.00,
			},
			prepare: func(f *fields) {

				f.operationsType.EXPECT().GetByID(gomock.Any(), 4).Times(1).Return(modelOperaTionsType.OperationType{
					OperationTypeID: 4,
					Description:     "PAGAMENTO",
					Operation:       1,
				}, nil)

				f.accounts.EXPECT().GetByID(gomock.Any(), "id").Times(1).Return(modelAccounts.Account{ID: "id"}, nil)

				f.transactions.EXPECT().Create(gomock.Any(), modelTransactions.MakeTransaction{
					AccountID:       "id",
					Amount:          60.00,
					OperationTypeID: 4,
				}).Times(1).Return(modelTransactions.Transaction{
					TransactionID:   "transaction_id",
					AccountID:       "id",
					Amount:          60.00,
					OperationTypeID: 4,
					Balance:         60,
				}, nil)

				f.wg.Add(3)

				f.transactions.EXPECT().GetToDischargeByAccountID(gomock.Any(), "id").Times(1).Return([]modelTransactions.Transaction{
					{
						TransactionID:   "1",
						AccountID:       "id",
						OperationTypeID: 1,
						Amount:          -60,
						Balance:         -60,
					},
				}, nil).Do(func(arg0, arg1 interface{}) {
					f.wg.Done()
				})

				f.transactions.EXPECT().UpdateBalance(gomock.Any(), modelTransactions.Transaction{
					TransactionID:   "1",
					AccountID:       "id",
					OperationTypeID: 1,
					Amount:          -60,
					Balance:         0,
				}).Times(1).Return(fmt.Errorf("any")).Do(func(arg0, arg1 interface{}) {
					f.wg.Done()
				})

				f.transactions.EXPECT().UpdateBalance(gomock.Any(), modelTransactions.Transaction{
					TransactionID:   "transaction_id",
					AccountID:       "id",
					OperationTypeID: 4,
					Amount:          60,
					Balance:         0,
				}).Times(1).Return(fmt.Errorf("any")).Do(func(arg0, arg1 interface{}) {
					f.wg.Done()
				})

			},
		},

		"should be able to make a new transaction with error to get in transactions": {
			input: modelTransactions.MakeTransaction{
				AccountID:       "id",
				OperationTypeID: 4,
				Amount:          60.00,
			},
			prepare: func(f *fields) {

				f.operationsType.EXPECT().GetByID(gomock.Any(), 4).Times(1).Return(modelOperaTionsType.OperationType{
					OperationTypeID: 4,
					Description:     "PAGAMENTO",
					Operation:       1,
				}, nil)

				f.accounts.EXPECT().GetByID(gomock.Any(), "id").Times(1).Return(modelAccounts.Account{ID: "id"}, nil)

				f.transactions.EXPECT().Create(gomock.Any(), modelTransactions.MakeTransaction{
					AccountID:       "id",
					Amount:          60.00,
					OperationTypeID: 4,
				}).Times(1).Return(modelTransactions.Transaction{
					TransactionID:   "transaction_id",
					AccountID:       "id",
					Amount:          60.00,
					OperationTypeID: 4,
					Balance:         60,
				}, nil)

				f.wg.Add(1)

				f.transactions.EXPECT().GetToDischargeByAccountID(gomock.Any(), "id").Times(1).Return(nil, fmt.Errorf("any")).Do(func(arg0, arg1 interface{}) {
					f.wg.Done()
				})

			},
		},

		"should be able to make a new transaction with dischard with empty transactions": {
			input: modelTransactions.MakeTransaction{
				AccountID:       "id",
				OperationTypeID: 4,
				Amount:          60.00,
			},
			prepare: func(f *fields) {

				f.operationsType.EXPECT().GetByID(gomock.Any(), 4).Times(1).Return(modelOperaTionsType.OperationType{
					OperationTypeID: 4,
					Description:     "PAGAMENTO",
					Operation:       1,
				}, nil)

				f.accounts.EXPECT().GetByID(gomock.Any(), "id").Times(1).Return(modelAccounts.Account{ID: "id"}, nil)

				f.transactions.EXPECT().Create(gomock.Any(), modelTransactions.MakeTransaction{
					AccountID:       "id",
					Amount:          60.00,
					OperationTypeID: 4,
				}).Times(1).Return(modelTransactions.Transaction{
					TransactionID:   "transaction_id",
					AccountID:       "id",
					Amount:          60.00,
					OperationTypeID: 4,
					Balance:         60,
				}, nil)

				f.wg.Add(1)

				f.transactions.EXPECT().GetToDischargeByAccountID(gomock.Any(), "id").Times(1).Return([]modelTransactions.Transaction{}, nil).Do(func(arg0, arg1 interface{}) {
					f.wg.Done()
				})
			},
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {

			ctrl := gomock.NewController(t)

			accountsMock := mocksStore.NewMockIAccounts(ctrl)
			transactionsMock := mocksStore.NewMockITransactions(ctrl)
			operationsTypeMock := mocksStore.NewMockIOperationsType(ctrl)
			var wg sync.WaitGroup

			tt.prepare(&fields{
				accounts:       accountsMock,
				transactions:   transactionsMock,
				operationsType: operationsTypeMock,
				wg:             &wg,
			})

			a := New(Options{
				Store: store.Store{
					Accounts:       accountsMock,
					Transactions:   transactionsMock,
					OperationsType: operationsTypeMock,
				},
				Log: logrus.New(),
			})

			err := a.Make(context.Background(), tt.input)
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf(`Expected err: "%s" got "%s"`, tt.err, err)
			}

			wg.Wait()

		})
	}
}
