package aws

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"go-import-from-s3/internal"
	"log"
	"strings"
	"time"
)

type DynamoDbClient struct {
	svc *dynamodb.DynamoDB
	cfg *internal.Config
}

func NewDynamoDbClient(c *internal.Config) *DynamoDbClient {
	return &DynamoDbClient{dynamodb.New(c.Session()), c}
}

func (s DynamoDbClient) Import() error {
	err := s.PrepareForImport()
	if err != nil {
		return err
	}

	importTableOutput, errImportTable := s.svc.ImportTable(s.getImportTableInput())
	if errImportTable != nil {
		return errImportTable
	}

	describeImportOutput := s.waitImportTable(importTableOutput)

	switch *describeImportOutput.ImportTableDescription.ImportStatus {
	case dynamodb.ImportStatusCompleted:
		log.Println("Importação do arquivo concluída")
	case dynamodb.ImportStatusCancelled, dynamodb.ImportStatusFailed:
		deleteTableError := s.deleteTable()
		if deleteTableError != nil {
			return deleteTableError
		}

		return errors.New(*describeImportOutput.ImportTableDescription.FailureMessage)
	}

	return nil
}

func (s DynamoDbClient) PrepareForImport() error {
	if describeTable, exists := s.tableExists(); exists {
		if *describeTable.Table.TableStatus != dynamodb.TableStatusActive {
			describeTable = s.waitFinalizationTableStatus()
		}

		if *describeTable.Table.TableStatus == dynamodb.TableStatusActive {
			err := s.deleteTable()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s DynamoDbClient) waitImportTable(importTable *dynamodb.ImportTableOutput) *dynamodb.DescribeImportOutput {
	for {
		describeImport, errDescribeImport := s.svc.DescribeImport(&dynamodb.DescribeImportInput{ImportArn: importTable.ImportTableDescription.ImportArn})
		if errDescribeImport != nil {
			log.Fatalln("Error > Import >", errDescribeImport)
		}

		switch *describeImport.ImportTableDescription.ImportStatus {
		case dynamodb.ImportStatusCompleted, dynamodb.ImportStatusCancelled, dynamodb.ImportStatusFailed:
			return describeImport
		default:
			log.Println("Aguardando a importação do arquivo...")
			time.Sleep(time.Second * 5)
		}
	}
}

func (s DynamoDbClient) waitFinalizationTableStatus() *dynamodb.DescribeTableOutput {
	var output *dynamodb.DescribeTableOutput
	var exists bool
	for {
		output, exists = s.tableExists()
		if exists == false ||
			*output.Table.TableStatus == dynamodb.TableStatusActive ||
			*output.Table.TableStatus == dynamodb.TableStatusArchived ||
			*output.Table.TableStatus == dynamodb.TableStatusInaccessibleEncryptionCredentials {
			break
		} else {
			log.Printf("A tabela %s está no status %s aguardando concluir...\n", s.cfg.Table(), *output.Table.TableStatus)
			time.Sleep(5 * time.Second)
		}
	}
	return output
}

func (s DynamoDbClient) deleteTable() error {
	log.Printf("Excluíndo a tabela %s", s.cfg.Table())
	output, err := s.svc.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(s.cfg.Table())})
	if err != nil {
		return fmt.Errorf("error > deleteTable > %s\n", err.Error())
	}

	if *output.TableDescription.TableStatus == dynamodb.TableStatusDeleting {
		for {
			resp, exists := s.tableExists()
			if exists && *resp.Table.TableStatus == dynamodb.TableStatusDeleting {
				log.Println("A tabela ainda está sendo excluída...")
				time.Sleep(5 * time.Second)
			} else {
				break
			}
		}
	}

	log.Printf("A tabela %s foi excluída com sucesso\n", s.cfg.Table())
	return nil
}

func (s DynamoDbClient) tableExists() (*dynamodb.DescribeTableOutput, bool) {
	output, err := s.svc.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(s.cfg.Table())})

	// Se houver um erro, verifica se é porque a tabela não existe
	if err != nil {
		awsError, ok := err.(awserr.Error)
		if ok && awsError.Code() == dynamodb.ErrCodeResourceNotFoundException {
			log.Printf("A tabela %s não foi encontrada\n", s.cfg.Table())
			return nil, false
		} else {
			log.Fatalln("Error > tableExists >", err)
		}
	}
	return output, true
}

func (s DynamoDbClient) EnableTimeToLive() error {
	if len(strings.TrimSpace(s.cfg.TtlName())) > 0 {
		_, err := s.svc.UpdateTimeToLive(&dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(s.cfg.Table()),
			TimeToLiveSpecification: &dynamodb.TimeToLiveSpecification{
				AttributeName: aws.String(s.cfg.TtlName()),
				Enabled:       aws.Bool(true),
			},
		})

		if err != nil {
			return err
		}

		for {
			output, errorDescribeTTL := s.svc.DescribeTimeToLive(&dynamodb.DescribeTimeToLiveInput{
				TableName: aws.String(s.cfg.Table()),
			})

			if errorDescribeTTL != nil {
				return errorDescribeTTL
			}

			if *output.TimeToLiveDescription.TimeToLiveStatus == dynamodb.TimeToLiveStatusEnabled {
				log.Println("TTL habilitado com sucesso na tabela", s.cfg.Table())
				break
			}

			log.Println("Aguardando habilitação do TTL na tabela", s.cfg.Table())
			time.Sleep(5 * time.Second)
		}
	}

	return nil
}

func (s DynamoDbClient) getImportTableInput() *dynamodb.ImportTableInput {
	importTableInput := &dynamodb.ImportTableInput{
		InputFormat: aws.String("CSV"),
		InputFormatOptions: &dynamodb.InputFormatOptions{
			Csv: &dynamodb.CsvOptions{
				Delimiter: aws.String(s.cfg.Delimiter()),
			},
		},
		S3BucketSource: &dynamodb.S3BucketSource{
			S3Bucket:    aws.String(s.cfg.Bucket()),
			S3KeyPrefix: aws.String(s.cfg.File()),
		},
		TableCreationParameters: &dynamodb.TableCreationParameters{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String(s.cfg.HashKey()),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String(s.cfg.RangeKey()),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String(s.cfg.HashKey()),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String(s.cfg.RangeKey()),
					KeyType:       aws.String("RANGE"),
				},
			},
			BillingMode: aws.String("PAY_PER_REQUEST"),
			TableName:   aws.String(s.cfg.Table()),
		},
	}
	return importTableInput
}
