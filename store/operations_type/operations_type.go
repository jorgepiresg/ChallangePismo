package operationsType

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	modelOperaTionsType "github.com/jorgepiresg/ChallangePismo/model/operations_type"
	"github.com/jorgepiresg/ChallangePismo/utils"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=$GOFILE -destination=../../mocks/store/operations_type_mock.go -package=mocksStore
type IOperationsType interface {
	GetByID(ctx context.Context, ID int) (modelOperaTionsType.OperationType, error)
}

type Options struct {
	DB    *sqlx.DB
	Log   *logrus.Logger
	Cache *redis.Client
}

type operationsType struct {
	db    *sqlx.DB
	log   *logrus.Logger
	cache *redis.Client
}

func New(opts Options) IOperationsType {
	return operationsType{
		db:    opts.DB,
		log:   opts.Log,
		cache: opts.Cache,
	}
}

func (ot operationsType) GetByID(ctx context.Context, ID int) (modelOperaTionsType.OperationType, error) {

	var operationsType modelOperaTionsType.OperationType

	cacheKey := fmt.Sprintf("operations_type_id_%d", ID)

	err := ot.getCache(ctx, cacheKey, &operationsType)
	if err == nil && !reflect.DeepEqual(operationsType, modelOperaTionsType.OperationType{}) {
		return operationsType, err
	}

	err = ot.db.GetContext(ctx, &operationsType, `SELECT operation_type_id, description, operation FROM operations_type where operation_type_id = $1`, ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			ot.log.WithField("operation_type_id_", ID).Error(err)
		}
		return operationsType, err
	}

	go ot.setCache(context.Background(), cacheKey, operationsType)

	return operationsType, nil
}

func (ot operationsType) setCache(ctx context.Context, key string, operationType modelOperaTionsType.OperationType) {
	err := ot.cache.Set(ctx, key, utils.ToJSON(operationType), 6*time.Hour).Err()
	if err != nil {
		ot.log.WithField("cache_key", key).Warning(err)
	}
}

func (ot operationsType) getCache(ctx context.Context, key string, operationType *modelOperaTionsType.OperationType) error {
	res, err := ot.cache.Get(ctx, key).Result()
	if err != nil {
		return nil
	}

	if err := utils.FromJson(res, operationType); err != nil {
		ot.log.WithField("cache_key", key).Error(err)
		return err
	}

	return nil
}
