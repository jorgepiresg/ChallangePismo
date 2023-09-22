package operationsType

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	modelOperaTionsType "github.com/jorgepiresg/ChallangePismo/model/operations_type"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestGetByID(t *testing.T) {

	type fields struct {
		sqlx sqlxmock.Sqlmock
	}

	tests := map[string]struct {
		input    int
		expected modelOperaTionsType.OperationType
		err      error
		prepare  func(f *fields)
	}{
		"success": {
			input: 1,
			prepare: func(f *fields) {
				rows := f.sqlx.NewRows([]string{"operation_type_id", "description", "operation"}).AddRow(1, "COMPRA A VISTA", -1)

				f.sqlx.ExpectQuery("SELECT operation_type_id, description, operation FROM operations_type").WithArgs(1).WillReturnRows(rows)
			},
			expected: modelOperaTionsType.OperationType{
				OperationTypeID: 1,
				Description:     "COMPRA A VISTA",
				Operation:       -1,
			},
		},
		"error": {
			input: 0,
			prepare: func(f *fields) {
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
