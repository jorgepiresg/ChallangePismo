package accounts

import (
	"context"

	"github.com/jmoiron/sqlx"
	modelAccounts "github.com/jorgepiresg/ChallangePismo/model/accounts"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/store/accounts_mock.go -package=mocksStore
type IAccounts interface {
	Create(ctx context.Context, account modelAccounts.Create) (modelAccounts.Account, error)
	GetByID(ctx context.Context, ID string) (modelAccounts.Account, error)
	GetByDocument(ctx context.Context, document string) (modelAccounts.Account, error)
}

type accounts struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) IAccounts {
	return accounts{
		db: db,
	}
}

func (a accounts) Create(ctx context.Context, create modelAccounts.Create) (modelAccounts.Account, error) {

	var account modelAccounts.Account

	rows, err := a.db.NamedQueryContext(ctx, `INSERT INTO accounts (document_number) VALUES (:document_number) RETURNING *`, create)
	if err != nil {
		return account, err
	}

	for rows.Next() {
		err = rows.StructScan(&account)
		if err != nil {
			return account, err
		}
	}

	return account, nil
}

func (a accounts) GetByID(ctx context.Context, ID string) (modelAccounts.Account, error) {

	var account modelAccounts.Account

	err := a.db.GetContext(ctx, &account, `SELECT account_id, document_number, created_at FROM accounts where account_id = $1`, ID)
	if err != nil {
		return account, err
	}
	return account, nil
}

func (a accounts) GetByDocument(ctx context.Context, document string) (modelAccounts.Account, error) {
	var account modelAccounts.Account

	err := a.db.GetContext(ctx, &account, `SELECT account_id, document_number, created_at FROM accounts where document_number = $1`, document)
	if err != nil {
		return account, err
	}
	return account, nil
}
