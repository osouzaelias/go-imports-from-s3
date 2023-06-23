package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"go-import-from-s3/internal"
	"log"
	"math"
	"path"
	"path/filepath"
	"time"
)

type S3Client struct {
	svc *s3.S3
	cfg *internal.Config
}

func NewS3Client(c *internal.Config) *S3Client {
	return &S3Client{
		s3.New(c.Session()),
		c,
	}
}

func (s S3Client) MoveToBackup() error {
	headOutput, _ := s.getHeadObject()

	objectSize := *headOutput.ContentLength
	partSize := int64(5 * math.Pow(2, 20))

	fileName := path.Join(s.cfg.Bucket(), s.cfg.File())
	backupFileName := fmt.Sprintf("%s/%s_%s", s.cfg.Backup(), time.Now().Format("20060102150405"), filepath.Base(s.cfg.File()))

	if objectSize < partSize {
		copyInput := &s3.CopyObjectInput{
			Bucket:     aws.String(s.cfg.Bucket()),
			CopySource: aws.String(fileName),
			Key:        aws.String(backupFileName),
		}

		if _, err := s.svc.CopyObject(copyInput); err != nil {
			return err
		}
	} else {

		var parts []*s3.CompletedPart

		createUploadInput := &s3.CreateMultipartUploadInput{
			Bucket: aws.String(s.cfg.Bucket()),
			Key:    aws.String(backupFileName),
		}

		createUploadOutput, createUploadError := s.svc.CreateMultipartUpload(createUploadInput)
		if createUploadError != nil {
			return createUploadError
		}

		uploadID := createUploadOutput.UploadId

		var partStart int64 = 0
		var partEnd = partStart + partSize - 1

		for i := int64(1); partStart < objectSize; i++ {

			partInput := &s3.UploadPartCopyInput{
				Bucket:          aws.String(s.cfg.Bucket()),
				CopySource:      aws.String(fileName),
				Key:             aws.String(backupFileName),
				CopySourceRange: aws.String(fmt.Sprintf("bytes=%d-%d", partStart, partEnd)),
				UploadId:        uploadID,
				PartNumber:      aws.Int64(i),
			}

			partOutput, partError := s.svc.UploadPartCopy(partInput)
			if partError != nil {
				return partError
			}

			parts = append(parts, &s3.CompletedPart{
				PartNumber: aws.Int64(i),
				ETag:       partOutput.CopyPartResult.ETag,
			})

			partStart += partSize

			if partEnd = partStart + partSize - 1; partEnd > objectSize {
				partEnd = objectSize - 1
			}

			bytesLeft := objectSize - partStart
			if bytesLeft < 0 {
				bytesLeft = 0
			}

			log.Printf("Upload da parte %d, restam ainda %d bytes", i, bytesLeft)
		}

		completeUploadInput := &s3.CompleteMultipartUploadInput{
			Bucket:          aws.String(s.cfg.Bucket()),
			Key:             aws.String(backupFileName),
			UploadId:        uploadID,
			MultipartUpload: &s3.CompletedMultipartUpload{Parts: parts},
		}

		_, completeUploadError := s.svc.CompleteMultipartUpload(completeUploadInput)

		if completeUploadError != nil {
			return completeUploadError
		}
	}

	log.Printf("Arquivo movido de %s para %s",
		fileName,
		backupFileName)

	return nil
}

func (s S3Client) DeleteFile() error {
	deleteInput := &s3.DeleteObjectInput{
		Bucket: aws.String(s.cfg.Bucket()),
		Key:    aws.String(s.cfg.File()),
	}

	_, err := s.svc.DeleteObject(deleteInput)
	if err != nil {
		return err
	}

	log.Printf("Arquivo %s excluído de %s", s.cfg.File(), s.cfg.Bucket())
	return nil
}

func (s S3Client) FileExists() bool {
	output, err := s.getHeadObject()

	if err != nil {
		log.Printf("Arquivo %s não encontrado", s.cfg.File())
		return false
	}

	log.Printf("Arquivo %s com %d bytes encontrado", s.cfg.File(), *output.ContentLength)

	return true
}

func (s S3Client) getHeadObject() (*s3.HeadObjectOutput, error) {
	return s.svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.cfg.Bucket()),
		Key:    aws.String(s.cfg.File()),
	})
}
