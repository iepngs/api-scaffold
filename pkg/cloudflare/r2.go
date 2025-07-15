package cloudflare

import (
	"bytes"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2Client Cloudflare R2 客户端
type R2Client struct {
	client *s3.Client
	bucket string
}

// NewR2Client 创建 R2 客户端
func NewR2Client(accessKey, secretKey, accountID, bucket string) *R2Client {
	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("auto"),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{URL: "https://" + accountID + ".r2.cloudflarestorage.com"}, nil
		})),
	)
	return &R2Client{client: s3.NewFromConfig(cfg), bucket: bucket}
}

// UploadFile 上传文件到 R2
func (r *R2Client) UploadFile(key string, body []byte) (string, error) {
	_, err := r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	})
	if err != nil {
		return "", err
	}
	return "https://" + r.bucket + "/" + key, nil
}

// DeleteFile 删除 R2 文件
func (r *R2Client) DeleteFile(key string) error {
	_, err := r.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	return err
}

// GetSignedURL 获取签名 URL
func (r *R2Client) GetSignedURL(key string, duration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(r.client)
	req, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(duration))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}
