package modelAccounts

import (
	"fmt"
	"strconv"
	"time"
)

type Account struct {
	ID             string    `json:"account_id,omitempty" db:"account_id"`
	DocumentNumber string    `json:"document_number,omitempty" db:"document_number"`
	CreatedAt      time.Time `json:"-" db:"created_at"`
}

type Create struct {
	DocumentNumber string `json:"document_number" db:"document_number"`
}

type CreateResponse struct {
	AccountID string `json:"account_id"`
}

func (c Create) Valid() error {

	if len(c.DocumentNumber) != 11 {
		return fmt.Errorf("document number invalid")

	}

	if _, err := strconv.Atoi(c.DocumentNumber); err != nil {
		return fmt.Errorf("document number invalid")
	}

	return nil
}
