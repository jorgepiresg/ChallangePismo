package transactions

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	modelTransactions "github.com/jorgepiresg/ChallangePismo/model/transactions"
	"github.com/sirupsen/logrus"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestCreate(t *testing.T) {
	type fields struct {
		sqlx sqlxmock.Sqlmock
	}

	tests := map[string]struct {
		input    modelTransactions.MakeTransaction
		expected modelTransactions.Transaction
		err      error
		prepare  func(f *fields)
	}{
		"should be able to insert transaction": {
			input: modelTransactions.MakeTransaction{
				AccountID:       "account_id",
				OperationTypeID: 1,
				Amount:          -10,
			},
			prepare: func(f *fields) {

				rows := f.sqlx.NewRows([]string{"transaction_id", "account_id", "operation_type_id", "amount", "event_date"}).AddRow("id", "account_id", 1, -10, time.Time{})

				f.sqlx.ExpectQuery("INSERT INTO transactions").WillReturnRows(rows)
			},
			expected: modelTransactions.Transaction{
				TransactionID:   "id",
				AccountID:       "account_id",
				OperationTypeID: 1,
				Amount:          -10,
			},
		},
		"should not be able to insert transaction with error at scan": {
			input: modelTransactions.MakeTransaction{
				AccountID:       "account_id",
				OperationTypeID: 1,
				Amount:          -10,
			},
			prepare: func(f *fields) {
				rows := f.sqlx.NewRows([]string{"id"}).AddRow("id")

				f.sqlx.ExpectQuery("INSERT INTO transactions").WillReturnRows(rows)
			},
			err: fmt.Errorf("missing destination name id in *modelTransactions.Transaction"),
		},
		"should not be able to insert transaction with error at sqlx": {
			input: modelTransactions.MakeTransaction{
				AccountID:       "account_id",
				OperationTypeID: 1,
				Amount:          -10,
			},
			prepare: func(f *fields) {
				f.sqlx.ExpectQuery("INSERT INTO transactions").WillReturnError(fmt.Errorf("any"))
			},
			err: fmt.Errorf("any"),
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {

			db, mock, err := sqlxmock.Newx()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			store := New(Options{
				DB:  db,
				Log: logrus.New(),
			})

			tt.prepare(&fields{
				sqlx: mock,
			})

			res, err := store.Create(context.Background(), tt.input)

			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf(`Expected err: "%s" got "%s"`, tt.err, err)
			}
			if !reflect.DeepEqual(res, tt.expected) {
				t.Errorf("Expected result %v got %v", tt.expected, res)
			}
		})
	}
}

func TestGetToDischargeByAccountID(t *testing.T) {

	type fields struct {
		sqlx sqlxmock.Sqlmock
	}

	tests := map[string]struct {
		input    string
		expected []modelTransactions.Transaction
		err      error
		prepare  func(f *fields)
	}{
		"should be able to get transactions to dischard by account id": {
			input: "1",
			prepare: func(f *fields) {

				rows := f.sqlx.NewRows([]string{"transaction_id", "account_id", "operation_type_id", "amount", "balance", "event_date"}).AddRow("1", "1", 1, -60, -60, time.Time{}).AddRow("2", "1", 1, -23.50, -23.50, time.Time{})

				f.sqlx.ExpectQuery("SELECT transaction_id ,account_id, operation_type_id, amount, balance, event_date FROM transactions").WithArgs("1").WillReturnRows(rows)

			},
			expected: []modelTransactions.Transaction{
				{
					TransactionID:   "1",
					AccountID:       "1",
					OperationTypeID: 1,
					Amount:          -60,
					Balance:         -60,
					EventDate:       time.Time{},
				},
				{
					TransactionID:   "2",
					AccountID:       "1",
					OperationTypeID: 1,
					Amount:          -23.50,
					Balance:         -23.50,
					EventDate:       time.Time{},
				},
			},
		},

		"should not be able to get transactions to dischard by account id with error": {
			input: "1",
			prepare: func(f *fields) {

				f.sqlx.ExpectQuery("SELECT transaction_id ,account_id, operation_type_id, amount, balance, event_date FROM transactions").WithArgs("1").WillReturnError(fmt.Errorf("any"))

			},
			err: fmt.Errorf("any"),
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {

			db, mock, err := sqlxmock.Newx()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			store := New(Options{
				DB:  db,
				Log: logrus.New(),
			})

			tt.prepare(&fields{
				sqlx: mock,
			})

			res, err := store.GetToDischargeByAccountID(context.Background(), tt.input)

			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf(`Expected err: "%s" got "%s"`, tt.err, err)
			}
			if !reflect.DeepEqual(res, tt.expected) {
				t.Errorf("Expected result %v got %v", tt.expected, res)
			}
		})
	}
}

func TestUpdateBalance(t *testing.T) {

	type fields struct {
		sqlx sqlxmock.Sqlmock
	}

	tests := map[string]struct {
		input   modelTransactions.Transaction
		err     error
		prepare func(f *fields)
	}{
		"should be able to update balance": {
			input: modelTransactions.Transaction{
				TransactionID:   "1",
				AccountID:       "1",
				OperationTypeID: 1,
				Amount:          -60,
				Balance:         0,
				EventDate:       time.Time{},
			},
			prepare: func(f *fields) {

				f.sqlx.ExpectExec("UPDATE transactions SET balance ").WithArgs(float64(0), "1").WillReturnResult(sqlxmock.NewResult(1, 1))

			},
		},

		"should not be able to update balance with error any": {
			input: modelTransactions.Transaction{
				TransactionID:   "1",
				AccountID:       "1",
				OperationTypeID: 1,
				Amount:          -60,
				Balance:         0,
				EventDate:       time.Time{},
			},
			err: fmt.Errorf("any"),
			prepare: func(f *fields) {

				f.sqlx.ExpectExec("UPDATE transactions SET balance ").WithArgs(float64(0), "1").WillReturnError(fmt.Errorf("any"))

			},
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {

			db, mock, err := sqlxmock.Newx()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			store := New(Options{
				DB:  db,
				Log: logrus.New(),
			})

			tt.prepare(&fields{
				sqlx: mock,
			})

			err = store.UpdateBalance(context.Background(), tt.input)

			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf(`Expected err: "%s" got "%s"`, tt.err, err)
			}

		})
	}
}
