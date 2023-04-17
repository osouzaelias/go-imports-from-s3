package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
)

type dynamodbHandler struct {
	db    *dynamodb.DynamoDB
	input *dynamodb.ImportTableInput
}

func NewDynamoDbHandler(c *config) (*dynamodbHandler, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(c.region),
	})

	if err != nil {
		log.Fatal(err)
	}

	svc := dynamodb.New(sess)

	return &dynamodbHandler{
		db:    svc,
		input: newImportTableInput(c),
	}, nil
}

func (d dynamodbHandler) ImportTable(fileName string, result interface{}) error {
	d.input.S3BucketSource.S3KeyPrefix = aws.String(fileName)
	result, err := d.db.ImportTable(d.input)

	if err != nil {
		return err
	}

	return nil
}

func (d dynamodbHandler) DescribeImport(id string, result interface{}) error {
	result, err := d.db.DescribeImport(&dynamodb.DescribeImportInput{
		ImportArn: aws.String(id),
	})

	if err != nil {
		return err
	}

	return nil
}

func newImportTableInput(c *config) *dynamodb.ImportTableInput {
	return &dynamodb.ImportTableInput{
		InputFormat: aws.String("CSV"),
		InputFormatOptions: &dynamodb.InputFormatOptions{
			Csv: &dynamodb.CsvOptions{
				Delimiter: aws.String(c.delimiter),
			},
		},
		S3BucketSource: &dynamodb.S3BucketSource{
			S3Bucket: aws.String(c.bucket),
		},
		TableCreationParameters: &dynamodb.TableCreationParameters{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String(c.hashKey),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("name"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String(c.hashKey),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String(c.rangeKey),
					KeyType:       aws.String("RANGE"),
				},
			},
			BillingMode: aws.String("PAY_PER_REQUEST"),
			TableName:   aws.String(c.table),
		},
	}
}
