package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
)

type S3Client struct {
	svc s3.S3
	cfg config
}

func NewS3Client() *S3Client {
	c := newConfigS3()
	return &S3Client{
		svc: *s3.New(&c.session),
		cfg: *c,
	}
}

func (s S3Client) MoveToBackup() {
	copyInput := &s3.CopyObjectInput{
		Bucket:     aws.String(s.cfg.bucket),
		CopySource: aws.String(s.cfg.bucket + "/" + s.cfg.file),
		Key:        aws.String(s.cfg.backup + s.cfg.file),
	}

	_, err := s.svc.CopyObject(copyInput)
	if err != nil {
		log.Fatalln("Error > MoveToBackup >", err)
	} else {
		log.Printf("Arquivo movido de %s para %s",
			aws.StringValue(copyInput.CopySource),
			aws.StringValue(copyInput.Key))
	}

	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(s.cfg.bucket),
		Key:    aws.String(s.cfg.file),
	}

	_, err = s.svc.DeleteObject(deleteInput)
	if err != nil {
		log.Fatalln("Error > MoveToBackup >", err)
	} else {
		log.Printf("Arquivo excluído de %s", aws.StringValue(copyInput.CopySource))
	}
}

func (s S3Client) FileExists() bool {
	output, err := s.svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.cfg.bucket),
		Key:    aws.String(s.cfg.file),
	})

	if err != nil {
		log.Println("Arquivo não encontrrado:", err)
		return false
	}

	log.Println("Arquivo encontrado, tamanho:", *output.ContentLength)

	return true
}
