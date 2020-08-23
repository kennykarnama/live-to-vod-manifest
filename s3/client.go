package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	s3cli "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
)

var (
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	svc        *s3cli.S3
)

func init() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	}))
	uploader = s3manager.NewUploader(sess)
	downloader = s3manager.NewDownloader(sess)
	svc = s3cli.New(sess)
}

func ListObjectsByPrefix(ctx context.Context, bucket, prefix string, delimiter string) (*s3cli.ListObjectsV2Output, error) {
	input := &s3cli.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}
	if delimiter != "" {
		input.Delimiter = aws.String(delimiter)
	}
	results, err := svc.ListObjectsV2(input)
	if err != nil {
		return nil, fmt.Errorf("action=ListObjectsByPrefix bucket=%v prefix=%v err=%v", bucket, prefix, err)
	}
	return results, nil
}

func GetWithContext(ctx context.Context, bucket, key string) (io.Reader, error) {
	input := &s3cli.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	result, err := svc.GetObjectWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func ListObjectPagesWithContext(ctx context.Context, bucket, prefix string, pageSize int64) ([]*s3cli.Object, error) {
	input := &s3cli.ListObjectsInput{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(pageSize),
		Prefix:  aws.String(prefix),
	}
	pageNum := 0
	var objectValues []*s3cli.Object
	err := svc.ListObjectsPagesWithContext(ctx, input, func(output *s3cli.ListObjectsOutput, lastPage bool) bool {
		pageNum++
		objectValues = append(objectValues, output.Contents...)
		if lastPage {
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	return objectValues, nil
}
