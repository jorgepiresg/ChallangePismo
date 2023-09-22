package accounts

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	modelAccounts "github.com/jorgepiresg/ChallangePismo/model/accounts"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestCreate(t *testing.T) {
	type fields struct {
		sqlx sqlxmock.Sqlmock
	}

	tests := map[string]struct {
		input    modelAccounts.Create
		expected modelAccounts.Account
		err      error
		prepare  func(f *fields)
	}{
		"success": {
			input: modelAccounts.Create{
				DocumentNumber: "111111111111",
			},
			prepare: func(f *fields) {
				rows := f.sqlx.NewRows([]string{"account_id", "document_number", "created_at"}).AddRow("id", "111111111111", time.Time{})

				f.sqlx.ExpectQuery("INSERT INTO accounts").WillReturnRows(rows)
			},
			expected: modelAccounts.Account{
				ID:             "id",
				DocumentNumber: "111111111111",
			},
		},
		"error scan": {
			input: modelAccounts.Create{
				DocumentNumber: "111111111111",
			},
			prepare: func(f *fields) {
				rows := f.sqlx.NewRows([]string{"id"}).AddRow("id")

				f.sqlx.ExpectQuery("INSERT INTO accounts").WillReturnRows(rows)
			},
			err: fmt.Errorf("missing destination name id in *modelAccounts.Account"),
		},
		"error": {
			input: modelAccounts.Create{
				DocumentNumber: "111111111111",
			},
			prepare: func(f *fields) {
				f.sqlx.ExpectQuery("INSERT INTO accounts").WillReturnError(fmt.Errorf("any"))
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

func TestGetByID(t *testing.T) {

	type fields struct {
		sqlx sqlxmock.Sqlmock
	}

	tests := map[string]struct {
		input    string
		expected modelAccounts.Account
		err      error
		prepare  func(f *fields)
	}{
		"success": {
			input: "id",
			prepare: func(f *fields) {
				rows := f.sqlx.NewRows([]string{"account_id", "document_number", "created_at"}).AddRow("id", "11111111111", time.Time{})

				f.sqlx.ExpectQuery("SELECT account_id, document_number, created_at FROM accounts").WithArgs("id").WillReturnRows(rows)
			},

			expected: modelAccounts.Account{
				ID:             "id",
				DocumentNumber: "11111111111",
			},
		},
		"error": {
			input: "invalid_id",
			prepare: func(f *fields) {
				f.sqlx.ExpectQuery("SELECT account_id, document_number, created_at FROM accounts").WithArgs("invalid_id").WillReturnError(fmt.Errorf("any"))
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

			res, err := store.GetByID(context.Background(), tt.input)

			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf(`Expected err: "%s" got "%s"`, tt.err, err)
			}
			if !reflect.DeepEqual(res, tt.expected) {
				t.Errorf("Expected result %v got %v", tt.expected, res)
			}
		})
	}
}

func TestGetByDocument(t *testing.T) {

	type fields struct {
		sqlx sqlxmock.Sqlmock
	}

	tests := map[string]struct {
		input    string
		expected modelAccounts.Account
		err      error
		prepare  func(f *fields)
	}{
		"success": {
			input: "11111111111",
			prepare: func(f *fields) {
				rows := f.sqlx.NewRows([]string{"account_id", "document_number", "created_at"}).AddRow("id", "11111111111", time.Time{})

				f.sqlx.ExpectQuery("SELECT account_id, document_number, created_at FROM accounts").WithArgs("11111111111").WillReturnRows(rows)
			},

			expected: modelAccounts.Account{
				ID:             "id",
				DocumentNumber: "11111111111",
			},
		},
		"error": {
			input: "11111111111",
			prepare: func(f *fields) {
				f.sqlx.ExpectQuery("SELECT account_id, document_number, created_at FROM accounts").WithArgs("11111111111").WillReturnError(fmt.Errorf("any"))
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

			res, err := store.GetByDocument(context.Background(), tt.input)

			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf(`Expected err: "%s" got "%s"`, tt.err, err)
			}
			if !reflect.DeepEqual(res, tt.expected) {
				t.Errorf("Expected result %v got %v", tt.expected, res)
			}
		})
	}
}
