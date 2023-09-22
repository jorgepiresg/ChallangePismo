package transactions

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	modelTransactions "github.com/jorgepiresg/ChallangePismo/model/transactions"
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
		"success": {
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
		"error scan": {
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
		"error": {
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

			store := New(db)

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
