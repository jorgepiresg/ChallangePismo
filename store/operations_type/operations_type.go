package operationsType

import (
	"context"

	"github.com/jmoiron/sqlx"
	modelOperaTionsType "github.com/jorgepiresg/ChallangePismo/model/operations_type"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/store/operations_type_mock.go -package=mocksStore
type IOperationsType interface {
	GetByID(ctx context.Context, ID int) (modelOperaTionsType.OperationType, error)
}

type operationsType struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) IOperationsType {
	return operationsType{
		db: db,
	}
}

func (ot operationsType) GetByID(ctx context.Context, ID int) (modelOperaTionsType.OperationType, error) {

	var operationsType modelOperaTionsType.OperationType

	err := ot.db.GetContext(ctx, &operationsType, `SELECT operation_type_id, description, operation FROM operations_type where operation_type_id = $1`, ID)
	if err != nil {
		return operationsType, err
	}
	return operationsType, nil
}
