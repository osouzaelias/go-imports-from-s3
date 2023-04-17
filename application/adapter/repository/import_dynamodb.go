package repository

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"go-import-from-s3/domain"
)

type ImportTableNoSQL struct {
	db NoSQL
}

func NewImportTableNoSQL(db NoSQL) ImportTableNoSQL {
	return ImportTableNoSQL{
		db: db,
	}
}

func (i ImportTableNoSQL) Import(table domain.ImportTable) (domain.ImportTable, error) {
	output := dynamodb.ImportTableOutput{}
	if err := i.db.ImportTable(table.FileName, output); err != nil {
		switch err.Error() {
		case dynamodb.ErrCodeResourceInUseException:
			return domain.ImportTable{}, errors.Wrap(domain.ErrResourceInUse, "error import table")
		case dynamodb.ErrCodeLimitExceededException:
			return domain.ImportTable{}, errors.Wrap(domain.ErrLimitExceeded, "error import table")
		case dynamodb.ErrCodeImportConflictException:
			return domain.ImportTable{}, errors.Wrap(domain.ErrImportConflict, "error import table")
		default:
			return domain.ImportTable{}, errors.Wrap(err, "error import table")
		}
	}

	return domain.ImportTable{
		ID:       *output.ImportTableDescription.ImportArn,
		FileName: *output.ImportTableDescription.S3BucketSource.S3KeyPrefix,
		Status:   *output.ImportTableDescription.ImportStatus,
	}, nil
}

func (i ImportTableNoSQL) Describe(table domain.ImportTable) (domain.ImportTable, error) {
	output := dynamodb.DescribeImportOutput{}
	if err := i.db.DescribeImport(table.ID, output); err != nil {
		switch err.Error() {
		case dynamodb.ErrCodeImportNotFoundException:
			return domain.ImportTable{}, errors.Wrap(domain.ErrImportNotFound, "error describe import")
		default:
			return domain.ImportTable{}, errors.Wrap(err, "error describe import")
		}
	}

	return domain.ImportTable{
		ID:       *output.ImportTableDescription.ImportArn,
		FileName: *output.ImportTableDescription.S3BucketSource.S3KeyPrefix,
		Status:   *output.ImportTableDescription.ImportStatus,
	}, nil
}
