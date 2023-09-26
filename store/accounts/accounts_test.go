package accounts

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	modelAccounts "github.com/jorgepiresg/ChallangePismo/model/accounts"
	"github.com/jorgepiresg/ChallangePismo/utils"
	"github.com/sirupsen/logrus"
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
		"should be able to insert account": {
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
		"should not be able to insert account with error at scan": {
			input: modelAccounts.Create{
				DocumentNumber: "111111111111",
			},
			prepare: func(f *fields) {
				rows := f.sqlx.NewRows([]string{"id"}).AddRow("id")

				f.sqlx.ExpectQuery("INSERT INTO accounts").WillReturnRows(rows)
			},
			err: fmt.Errorf("missing destination name id in *modelAccounts.Account"),
		},
		"should not be able to insert account with error at sqlx": {
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

func TestGetByID(t *testing.T) {

	type fields struct {
		sqlx  sqlxmock.Sqlmock
		redis redismock.ClientMock
	}

	tests := map[string]struct {
		input    string
		expected modelAccounts.Account
		err      error
		prepare  func(f *fields)
	}{
		"should be able to get account by id": {
			input: "id",
			prepare: func(f *fields) {

				f.redis.ExpectGet("account_id_id").RedisNil()

				rows := f.sqlx.NewRows([]string{"account_id", "document_number", "created_at"}).AddRow("id", "11111111111", time.Time{})

				f.sqlx.ExpectQuery("SELECT account_id, document_number, created_at FROM accounts").WithArgs("id").WillReturnRows(rows)

				f.redis.ExpectSet("account_id_id", utils.ToJSON(modelAccounts.Account{
					ID:             "id",
					DocumentNumber: "11111111111",
				}), time.Minute).SetVal("")

			},
			expected: modelAccounts.Account{
				ID:             "id",
				DocumentNumber: "11111111111",
			},
		},

		"should be able to get account by id with error to save at cache": {
			input: "id",
			prepare: func(f *fields) {

				f.redis.ExpectGet("account_id_id").RedisNil()

				rows := f.sqlx.NewRows([]string{"account_id", "document_number", "created_at"}).AddRow("id", "11111111111", time.Time{})

				f.sqlx.ExpectQuery("SELECT account_id, document_number, created_at FROM accounts").WithArgs("id").WillReturnRows(rows)

				f.redis.ExpectSet("account_id_id", utils.ToJSON(modelAccounts.Account{
					ID:             "id",
					DocumentNumber: "11111111111",
				}), time.Minute).SetErr(fmt.Errorf("any"))

			},
			expected: modelAccounts.Account{
				ID:             "id",
				DocumentNumber: "11111111111",
			},
		},

		"should be able to get account by id with error to unmarshal from cache": {
			input: "id",
			prepare: func(f *fields) {

				f.redis.ExpectGet("account_id_id").SetVal(`A`)

				rows := f.sqlx.NewRows([]string{"account_id", "document_number", "created_at"}).AddRow("id", "11111111111", time.Time{})

				f.sqlx.ExpectQuery("SELECT account_id, document_number, created_at FROM accounts").WithArgs("id").WillReturnRows(rows)

				f.redis.ExpectSet("account_id_id", utils.ToJSON(modelAccounts.Account{
					ID:             "id",
					DocumentNumber: "11111111111",
				}), time.Minute).SetErr(fmt.Errorf("any"))

			},
			expected: modelAccounts.Account{
				ID:             "id",
				DocumentNumber: "11111111111",
			},
		},

		"should be able to get account by id in cache": {
			input: "id",
			prepare: func(f *fields) {

				f.redis.ExpectGet("account_id_id").SetVal(`{"account_id":"id", "document_number":"11111111111"}`)

			},
			expected: modelAccounts.Account{
				ID:             "id",
				DocumentNumber: "11111111111",
			},
		},

		"should not be able to get account by id with error at sqlx": {
			input: "invalid_id",
			prepare: func(f *fields) {

				f.redis.ExpectGet("account_id_id").RedisNil()

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

			cacheDB, cacheMock := redismock.NewClientMock()

			store := New(Options{
				DB:    db,
				Log:   logrus.New(),
				Cache: cacheDB,
			})

			tt.prepare(&fields{
				sqlx:  mock,
				redis: cacheMock,
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
		sqlx  sqlxmock.Sqlmock
		redis redismock.ClientMock
	}

	tests := map[string]struct {
		input    string
		expected modelAccounts.Account
		err      error
		prepare  func(f *fields)
	}{
		"should be able to get account by document": {
			input: "11111111111",
			prepare: func(f *fields) {

				f.redis.ExpectGet("account_document_11111111111").RedisNil()

				rows := f.sqlx.NewRows([]string{"account_id", "document_number", "created_at"}).AddRow("id", "11111111111", time.Time{})

				f.sqlx.ExpectQuery("SELECT account_id, document_number, created_at FROM accounts").WithArgs("11111111111").WillReturnRows(rows)

				f.redis.ExpectSet("account_document_11111111111", utils.ToJSON(modelAccounts.Account{
					ID:             "id",
					DocumentNumber: "11111111111",
				}), time.Minute).SetVal("")
			},

			expected: modelAccounts.Account{
				ID:             "id",
				DocumentNumber: "11111111111",
			},
		},
		"should be able to get account by document with error to save at cache": {
			input: "11111111111",
			prepare: func(f *fields) {

				f.redis.ExpectGet("account_document_11111111111").RedisNil()

				rows := f.sqlx.NewRows([]string{"account_id", "document_number", "created_at"}).AddRow("id", "11111111111", time.Time{})

				f.sqlx.ExpectQuery("SELECT account_id, document_number, created_at FROM accounts").WithArgs("11111111111").WillReturnRows(rows)

				f.redis.ExpectSet("account_document_11111111111", utils.ToJSON(modelAccounts.Account{
					ID:             "id",
					DocumentNumber: "11111111111",
				}), time.Minute).SetErr(fmt.Errorf("any"))

			},
			expected: modelAccounts.Account{
				ID:             "id",
				DocumentNumber: "11111111111",
			},
		},

		"should be able to get account by document with error to unmarshal from cache": {
			input: "11111111111",
			prepare: func(f *fields) {

				f.redis.ExpectGet("account_document_11111111111").SetVal(`A`)

				rows := f.sqlx.NewRows([]string{"account_id", "document_number", "created_at"}).AddRow("id", "11111111111", time.Time{})

				f.sqlx.ExpectQuery("SELECT account_id, document_number, created_at FROM accounts").WithArgs("11111111111").WillReturnRows(rows)

				f.redis.ExpectSet("account_document_11111111111", utils.ToJSON(modelAccounts.Account{
					ID:             "id",
					DocumentNumber: "11111111111",
				}), time.Minute).SetErr(fmt.Errorf("any"))

			},
			expected: modelAccounts.Account{
				ID:             "id",
				DocumentNumber: "11111111111",
			},
		},

		"should be able to get account by document in cache": {
			input: "11111111111",
			prepare: func(f *fields) {

				f.redis.ExpectGet("account_document_11111111111").SetVal(`{"account_id":"id", "document_number":"11111111111"}`)

			},
			expected: modelAccounts.Account{
				ID:             "id",
				DocumentNumber: "11111111111",
			},
		},
		"should not be able to get account by document with error at sqlx": {
			input: "11111111111",
			prepare: func(f *fields) {

				f.redis.ExpectGet("account_document_11111111111").RedisNil()

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
			cacheDB, cacheMock := redismock.NewClientMock()

			store := New(Options{
				DB:    db,
				Log:   logrus.New(),
				Cache: cacheDB,
			})

			tt.prepare(&fields{
				sqlx:  mock,
				redis: cacheMock,
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
