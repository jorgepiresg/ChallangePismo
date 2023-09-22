package accounts

import (
	"context"
	"fmt"

	modelAccounts "github.com/jorgepiresg/ChallangePismo/model/accounts"
	"github.com/jorgepiresg/ChallangePismo/store"
	"github.com/jorgepiresg/ChallangePismo/utils"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/app/accounts_mock.go -package=mocksApp
type IAccounts interface {
	Create(ctx context.Context, account modelAccounts.Create) (modelAccounts.Account, error)
	GetByAccountID(ctx context.Context, AccountID string) (modelAccounts.Account, error)
}

type Options struct {
	Store store.Store
}

type account struct {
	store store.Store
}

func New(opts Options) IAccounts {
	return account{
		store: opts.Store,
	}
}

func (a account) Create(ctx context.Context, account modelAccounts.Create) (modelAccounts.Account, error) {

	var emptyAccount modelAccounts.Account

	account.DocumentNumber = utils.CleanDocument(account.DocumentNumber)

	if err := account.Valid(); err != nil {
		return emptyAccount, err
	}

	if _, err := a.store.Accounts.GetByDocument(ctx, account.DocumentNumber); err == nil {
		return emptyAccount, fmt.Errorf("user alredy exist")
	}

	return a.store.Accounts.Create(ctx, account)
}

func (a account) GetByAccountID(ctx context.Context, AccountID string) (modelAccounts.Account, error) {
	return a.store.Accounts.GetByID(ctx, AccountID)
}
