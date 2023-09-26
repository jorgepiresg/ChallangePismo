package operationsType

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	modelOperaTionsType "github.com/jorgepiresg/ChallangePismo/model/operations_type"
	"github.com/jorgepiresg/ChallangePismo/utils"
	"github.com/sirupsen/logrus"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestGetByID(t *testing.T) {

	type fields struct {
		sqlx  sqlxmock.Sqlmock
		redis redismock.ClientMock
	}

	tests := map[string]struct {
		input    int
		expected modelOperaTionsType.OperationType
		err      error
		prepare  func(f *fields)
	}{
		"should be able to get operation type by id": {
			input: 1,
			prepare: func(f *fields) {

				f.redis.ExpectGet("operations_type_id_1").RedisNil()

				rows := f.sqlx.NewRows([]string{"operation_type_id", "description", "operation"}).AddRow(1, "COMPRA A VISTA", -1)

				f.sqlx.ExpectQuery("SELECT operation_type_id, description, operation FROM operations_type").WithArgs(1).WillReturnRows(rows)

				f.redis.ExpectSet("operations_type_id_1", utils.ToJSON(modelOperaTionsType.OperationType{
					OperationTypeID: 1,
					Description:     "COMPRA A VISTA",
					Operation:       -1,
				}), 6*time.Hour).SetVal("")
			},
			expected: modelOperaTionsType.OperationType{
				OperationTypeID: 1,
				Description:     "COMPRA A VISTA",
				Operation:       -1,
			},
		},

		"should be able to get operation type by id with error to save at cache": {
			input: 1,
			prepare: func(f *fields) {

				f.redis.ExpectGet("operations_type_id_1").RedisNil()

				rows := f.sqlx.NewRows([]string{"operation_type_id", "description", "operation"}).AddRow(1, "COMPRA A VISTA", -1)

				f.sqlx.ExpectQuery("SELECT operation_type_id, description, operation FROM operations_type").WithArgs(1).WillReturnRows(rows)

				f.redis.ExpectSet("operations_type_id_1", utils.ToJSON(modelOperaTionsType.OperationType{
					OperationTypeID: 1,
					Description:     "COMPRA A VISTA",
					Operation:       -1,
				}), 6*time.Hour).SetErr(fmt.Errorf("any"))
			},
			expected: modelOperaTionsType.OperationType{
				OperationTypeID: 1,
				Description:     "COMPRA A VISTA",
				Operation:       -1,
			},
		},

		"should be able to get operation type by id with error to unmarshal from cache": {
			input: 1,
			prepare: func(f *fields) {

				f.redis.ExpectGet("operations_type_id_1").SetVal(`A`)

				rows := f.sqlx.NewRows([]string{"operation_type_id", "description", "operation"}).AddRow(1, "COMPRA A VISTA", -1)

				f.sqlx.ExpectQuery("SELECT operation_type_id, description, operation FROM operations_type").WithArgs(1).WillReturnRows(rows)

				f.redis.ExpectSet("operations_type_id_1", utils.ToJSON(modelOperaTionsType.OperationType{
					OperationTypeID: 1,
					Description:     "COMPRA A VISTA",
					Operation:       -1,
				}), 6*time.Hour).SetVal("")
			},
			expected: modelOperaTionsType.OperationType{
				OperationTypeID: 1,
				Description:     "COMPRA A VISTA",
				Operation:       -1,
			},
		},

		"should be able to get operation type by id in cache": {
			input: 1,
			prepare: func(f *fields) {

				f.redis.ExpectGet("operations_type_id_1").SetVal(`{"operation_type_id":1, "description":"COMPRA A VISTA", "operation":-1}`)

			},
			expected: modelOperaTionsType.OperationType{
				OperationTypeID: 1,
				Description:     "COMPRA A VISTA",
				Operation:       -1,
			},
		},

		"should not be able to get operation type by id with error at sqlx": {
			input: 0,
			prepare: func(f *fields) {
				f.redis.ExpectGet("operation_type_id_1").RedisNil()

				f.sqlx.ExpectQuery("SELECT operation_type_id, description, operation FROM operations_type").WillReturnError(fmt.Errorf("any"))
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
