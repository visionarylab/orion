package storage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-pg/pg/orm"

	"github.com/Syncano/orion/pkg/settings"
	"github.com/Syncano/orion/pkg/util"
)

type s3Storage struct {
	uploader *s3manager.Uploader
	client   *s3.S3
}

func newS3Storage() DataStorage {
	session := createS3Session(settings.Storage.AccessKeyID, settings.Storage.SecretAccessKey,
		settings.Storage.Region, settings.Storage.Endpoint)
	client := s3.New(session)
	uploader := s3manager.NewUploaderWithClient(client)

	return &s3Storage{
		uploader: uploader,
		client:   client,
	}
}

// Client returns s3 client.
func (s *s3Storage) Client() interface{} {
	return s.client
}

func createS3Session(accessKeyID, secretAccessKey, region, endpoint string) *session.Session {
	conf := aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	}
	creds := credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")
	sess, err := session.NewSession(conf.WithCredentials(creds))
	util.Must(err)
	return sess
}

func (s *s3Storage) SafeUpload(ctx context.Context, db orm.DB, bucket, key string, f io.Reader) error {
	AddDBRollbackHook(db, func() {
		util.Must(
			s.Delete(ctx, bucket, key),
		)
	})
	return s.Upload(ctx, bucket, key, f)
}

func (s *s3Storage) Upload(ctx context.Context, bucket, key string, f io.Reader) error {
	_, err := s.uploader.UploadWithContext(ctx,
		&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			ACL:    aws.String("public-read"),
			Body:   f,
		})
	return err
}

func (s *s3Storage) Delete(ctx context.Context, bucket, key string) error {
	_, err := s.client.DeleteObjectWithContext(
		ctx,
		&s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	return err
}
