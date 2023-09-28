package modelTransactions

import (
	"fmt"
	"strings"
	"time"
)

type Transaction struct {
	TransactionID   string    `db:"transaction_id"`
	AccountID       string    `db:"account_id"`
	OperationTypeID int       `db:"operation_type_id"`
	Amount          float64   `db:"amount"`
	Balance         float64   `db:"balance"`
	EventDate       time.Time `db:"event_date"`
}

type MakeTransaction struct {
	AccountID       string  `json:"account_id" db:"account_id"`
	OperationTypeID int     `json:"operation_type_id" db:"operation_type_id"`
	Amount          float64 `json:"amount" db:"amount"`
}

func (dt *MakeTransaction) SetOperationInAmount(operation int) error {

	dt.Amount = dt.Amount * float64(operation)

	return nil
}

func (dt *MakeTransaction) ValidateAmount() error {

	if dt.Amount == 0 {
		return fmt.Errorf("amount 0 is invalid")
	}

	amountString := fmt.Sprintf("%v", dt.Amount)
	split := strings.Split(amountString, ".")

	if len(split) > 1 && len(split[1]) != 2 {
		return fmt.Errorf("amount %v is invalid, use 2 decimals", dt.Amount)
	}

	return nil
}
