package internal

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

//const bucketName = "bucketelias"
//const s.cfg.file = "myFile0.csv"
//const backupFolder = "backup/"

type ServiceS3 struct {
	svc s3.S3
	cfg Config
}

func NewServiceS3(c Config) *ServiceS3 {
	return &ServiceS3{
		svc: *s3.New(&c.session),
		cfg: c,
	}
}

func (s ServiceS3) MoveToBackup() error {
	copyInput := &s3.CopyObjectInput{
		Bucket:     aws.String(s.cfg.bucket),
		CopySource: aws.String(s.cfg.bucket + "/" + s.cfg.file),
		Key:        aws.String(s.cfg.backup + s.cfg.file),
	}

	_, err := s.svc.CopyObject(copyInput)
	if err != nil {
		return err
	}

	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(s.cfg.bucket),
		Key:    aws.String(s.cfg.file),
	}
	_, err = s.svc.DeleteObject(deleteInput)
	return err
}

func (s ServiceS3) FileExists() bool {
	output, err := s.svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.cfg.bucket),
		Key:    aws.String(s.cfg.file),
	})

	if err != nil {
		fmt.Println("Arquivo não encontrrado, erro:", err)
		return false
	}

	fmt.Println("Arquivo encontrado, tamanho:", aws.Int64Value(output.ContentLength))

	return true
}
