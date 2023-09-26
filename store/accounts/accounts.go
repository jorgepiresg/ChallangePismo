package accounts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	modelAccounts "github.com/jorgepiresg/ChallangePismo/model/accounts"
	"github.com/jorgepiresg/ChallangePismo/utils"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/store/accounts_mock.go -package=mocksStore
type IAccounts interface {
	Create(ctx context.Context, account modelAccounts.Create) (modelAccounts.Account, error)
	GetByID(ctx context.Context, ID string) (modelAccounts.Account, error)
	GetByDocument(ctx context.Context, document string) (modelAccounts.Account, error)
}

type Options struct {
	DB    *sqlx.DB
	Log   *logrus.Logger
	Cache *redis.Client
}

type accounts struct {
	db    *sqlx.DB
	log   *logrus.Logger
	cache *redis.Client
}

func New(opts Options) IAccounts {
	return accounts{
		db:    opts.DB,
		log:   opts.Log,
		cache: opts.Cache,
	}
}

func (a accounts) Create(ctx context.Context, create modelAccounts.Create) (modelAccounts.Account, error) {

	var account modelAccounts.Account

	rows, err := a.db.NamedQueryContext(ctx, `INSERT INTO accounts (document_number) VALUES (:document_number) RETURNING *`, create)
	if err != nil {
		a.log.WithField("document", create.DocumentNumber).Error(err)
		return account, err
	}

	for rows.Next() {
		err = rows.StructScan(&account)
		if err != nil {
			a.log.WithField("document", create.DocumentNumber).Error(err)
			return account, err
		}
	}

	return account, nil
}

func (a accounts) GetByID(ctx context.Context, ID string) (modelAccounts.Account, error) {

	var account modelAccounts.Account
	cacheKey := fmt.Sprintf("account_id_%s", ID)

	err := a.getCache(ctx, cacheKey, &account)
	if err == nil && !reflect.DeepEqual(account, modelAccounts.Account{}) {
		return account, err
	}

	err = a.db.GetContext(ctx, &account, `SELECT account_id, document_number, created_at FROM accounts where account_id = $1`, ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			a.log.WithField("account_id", ID).Error(err)
		}
		return account, err
	}

	go a.setCache(context.Background(), cacheKey, account)

	return account, nil
}

func (a accounts) GetByDocument(ctx context.Context, document string) (modelAccounts.Account, error) {

	var account modelAccounts.Account

	cacheKey := fmt.Sprintf("account_document_%s", document)

	err := a.getCache(ctx, cacheKey, &account)
	if err == nil && !reflect.DeepEqual(account, modelAccounts.Account{}) {
		return account, err
	}

	err = a.db.GetContext(ctx, &account, `SELECT account_id, document_number, created_at FROM accounts where document_number = $1`, document)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			a.log.WithField("document", document).Error(err)
		}
		return account, err
	}

	go a.setCache(context.Background(), cacheKey, account)

	return account, nil
}

func (a accounts) setCache(ctx context.Context, key string, account modelAccounts.Account) {

	err := a.cache.Set(ctx, key, utils.ToJSON(account), 10*time.Minute).Err()
	if err != nil {
		a.log.WithField("cache_key", key).Warning(err)
	}
}

func (a accounts) getCache(ctx context.Context, key string, account *modelAccounts.Account) error {
	res, err := a.cache.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	if err := utils.FromJson(res, account); err != nil {
		a.log.WithField("cache_key", key).Error(err)
		return err
	}

	return nil
}
