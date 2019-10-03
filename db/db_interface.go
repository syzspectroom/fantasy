package db

import (
	"context"
)

// DbInterface database interface
type DbInterface interface {
	Insert(ctx context.Context, colName string, obj interface{}) error
	InsertMany(ctx context.Context, colName string, obj interface{}) error
	Query(ctx context.Context, query string, bindVars map[string]interface{}, resObj interface{}) (string, error)
	Update(ctx context.Context, colName string, key string, obj interface{}) error
}
