package aws

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"strings"
	"time"
)

type DynamoDbClient struct {
	svc *dynamodb.DynamoDB
	cfg *Config
}

func NewDynamoDbClient(c *Config) *DynamoDbClient {
	return &DynamoDbClient{dynamodb.New(c.session), c}
}

func (s DynamoDbClient) Import() error {
	s.prepareForImport()

	importTableOutput, errImportTable := s.svc.ImportTable(s.getImportTableInput())
	if errImportTable != nil {
		return errImportTable
	}

	describeImportOutput := s.waitImportTable(importTableOutput)

	switch *describeImportOutput.ImportTableDescription.ImportStatus {
	case dynamodb.ImportStatusCompleted:
		log.Println("Importação do arquivo concluída")
	case dynamodb.ImportStatusCancelled, dynamodb.ImportStatusFailed:
		s.deleteTable()
		return errors.New(*describeImportOutput.ImportTableDescription.FailureMessage)
	}

	return nil
}

func (s DynamoDbClient) prepareForImport() {
	if describeTable, exists := s.tableExists(); exists {
		if *describeTable.Table.TableStatus != dynamodb.TableStatusActive {
			describeTable = s.waitFinalizationTableStatus()
		}

		if *describeTable.Table.TableStatus == dynamodb.TableStatusActive {
			s.deleteTable()
		}
	}
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
			log.Printf("A tabela %s está no status %s aguardando concluir...\n", s.cfg.table, *output.Table.TableStatus)
			time.Sleep(5 * time.Second)
		}
	}
	return output
}

func (s DynamoDbClient) deleteTable() {
	log.Printf("Excluíndo a tabela %s", s.cfg.table)
	output, err := s.svc.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(s.cfg.table)})
	if err != nil {
		log.Fatalln("Error > deleteTable >", err)
	}

	if *output.TableDescription.TableStatus == dynamodb.TableStatusDeleting {
		for {
			resp, exists := s.tableExists()
			if exists && aws.StringValue(resp.Table.TableStatus) == dynamodb.TableStatusDeleting {
				log.Println("A tabela ainda está sendo excluída...")
				time.Sleep(5 * time.Second)
			} else {
				break
			}
		}
	}

	log.Printf("A tabela %s foi excluída com sucesso\n", s.cfg.table)
}

func (s DynamoDbClient) tableExists() (*dynamodb.DescribeTableOutput, bool) {
	output, err := s.svc.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(s.cfg.table)})

	// Se houver um erro, verifica se é porque a tabela não existe
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
			log.Printf("A tabela %s não foi encontrada\n", s.cfg.table)
			return nil, false
		} else {
			log.Fatalln("Error > tableExists >", err)
		}
	}
	return output, true
}

func (s DynamoDbClient) EnableTimeToLive() error {
	if len(strings.TrimSpace(s.cfg.ttlAttributeName)) > 0 {
		_, err := s.svc.UpdateTimeToLive(&dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(s.cfg.table),
			TimeToLiveSpecification: &dynamodb.TimeToLiveSpecification{
				AttributeName: aws.String(s.cfg.ttlAttributeName),
				Enabled:       aws.Bool(true),
			},
		})

		if err != nil {
			return err
		}

		for {
			output, errorDescribeTTL := s.svc.DescribeTimeToLive(&dynamodb.DescribeTimeToLiveInput{
				TableName: aws.String(s.cfg.table),
			})

			if errorDescribeTTL != nil {
				return errorDescribeTTL
			}

			if *output.TimeToLiveDescription.TimeToLiveStatus == dynamodb.TimeToLiveStatusEnabled {
				log.Println("TTL habilitado com sucesso na tabela", s.cfg.table)
				break
			}

			log.Println("Aguardando habilitação do TTL na tabela", s.cfg.table)
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
				Delimiter: aws.String(s.cfg.delimiter),
			},
		},
		S3BucketSource: &dynamodb.S3BucketSource{
			S3Bucket:    aws.String(s.cfg.bucket),
			S3KeyPrefix: aws.String(s.cfg.file),
		},
		TableCreationParameters: &dynamodb.TableCreationParameters{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String(s.cfg.hashKey),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String(s.cfg.rangeKey),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String(s.cfg.hashKey),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String(s.cfg.rangeKey),
					KeyType:       aws.String("RANGE"),
				},
			},
			BillingMode: aws.String("PAY_PER_REQUEST"),
			TableName:   aws.String(s.cfg.table),
		},
	}
	return importTableInput
}
