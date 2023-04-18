package internal

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"time"
)

//const tableName = "tb-import-from-s3-v6"
//const hashKey = "id"
//const rangeKey = "name"

type ServiceDynamoDb struct {
	svc dynamodb.DynamoDB
	cfg Config
}

func NewServiceDynamoDb(c Config) *ServiceDynamoDb {
	return &ServiceDynamoDb{
		svc: *dynamodb.New(&c.session),
		cfg: c,
	}
}

func (s ServiceDynamoDb) Import() error {
	if describeTable, exists := s.tableExists(); exists {
		if aws.StringValue(describeTable.Table.TableStatus) != dynamodb.TableStatusActive {
			describeTable = s.waitFinalizationTableStatus()
		}

		if aws.StringValue(describeTable.Table.TableStatus) == dynamodb.TableStatusActive {
			s.deleteTable()
		}
	}

	importTable, err := s.svc.ImportTable(s.getImportTableInput())
	if err != nil {
		return err
	}

	for {
		describeImport, err := s.svc.DescribeImport(&dynamodb.DescribeImportInput{ImportArn: importTable.ImportTableDescription.ImportArn})
		if err != nil {
			return err
		}

		switch *describeImport.ImportTableDescription.ImportStatus {
		case dynamodb.ImportStatusCompleted:
			return nil
		case dynamodb.ImportStatusCancelled, dynamodb.ImportStatusFailed:
			log.Println("A importação foi interrompida")
			s.deleteTable()
			log.Fatalln("Error > Import:", aws.StringValue(describeImport.ImportTableDescription.FailureMessage))
		default:
			log.Println("Aguardando a importação do arquivo...")
			time.Sleep(time.Second * 5)
		}
	}
}

func (s ServiceDynamoDb) waitFinalizationTableStatus() *dynamodb.DescribeTableOutput {
	var output *dynamodb.DescribeTableOutput
	var exists bool
	for {
		output, exists = s.tableExists()
		if exists == false ||
			aws.StringValue(output.Table.TableStatus) == dynamodb.TableStatusActive ||
			aws.StringValue(output.Table.TableStatus) == dynamodb.TableStatusArchived ||
			aws.StringValue(output.Table.TableStatus) == dynamodb.TableStatusInaccessibleEncryptionCredentials {
			break
		} else {
			log.Printf("A tabela %s está no status %s aguardando concluir...\n", s.cfg.table, aws.StringValue(output.Table.TableStatus))
			time.Sleep(5 * time.Second)
		}
	}
	return output
}

func (s ServiceDynamoDb) deleteTable() {
	log.Printf("Excluíndo a tabela %s", s.cfg.table)
	output, err := s.svc.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(s.cfg.table)})
	if err != nil {
		log.Fatalln("Error > deleteTable:", err)
	}

	if aws.StringValue(output.TableDescription.TableStatus) == dynamodb.TableStatusDeleting {
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

func (s ServiceDynamoDb) tableExists() (*dynamodb.DescribeTableOutput, bool) {
	output, err := s.svc.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(s.cfg.table)})

	// Se houver um erro, verifica se é porque a tabela não existe
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok && aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
			log.Printf("A tabela %s não foi encontrada\n", s.cfg.table)
			return nil, false
		} else {
			log.Fatalln("Error > tableExists:", err)
		}
	}
	return output, true
}

func (s ServiceDynamoDb) getImportTableInput() *dynamodb.ImportTableInput {
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
