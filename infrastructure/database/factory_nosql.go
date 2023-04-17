package database

import (
	"errors"
	"go-import-from-s3/application/adapter/repository"
)

var errInvalidNoSQLDatabaseInstance = errors.New("invalid nosql db instance")

const (
	InstanceDynamoDB int = iota
)

func NewDatabaseNoSQLFactory(instance int) (repository.NoSQL, error) {
	switch instance {
	case InstanceDynamoDB:
		return NewDynamoDbHandler(newConfigDynamoDB())
	default:
		return nil, errInvalidNoSQLDatabaseInstance
	}
}
