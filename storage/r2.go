package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"backend-crowdfunding/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type r2Uploader struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

func newR2Uploader(cfg config.StorageConfig) (*r2Uploader, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		// R2 ignores the region but the SDK still requires one.
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint())
	})

	return &r2Uploader{
		client:    client,
		bucket:    cfg.Bucket,
		publicURL: strings.TrimRight(cfg.PublicURL, "/"),
	}, nil
}

func (u *r2Uploader) Upload(ctx context.Context, key string, body io.Reader, size int64, contentType string) (string, error) {
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(u.bucket),
		Key:           aws.String(key),
		Body:          body,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", u.publicURL, key), nil
}
