package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	bucket = "amz-cdn-s3-testing"
	region = ""
)

type S3Client struct {
	Client *s3.Client
	Bucket string
	Region string
}

func main() {
	// create root context
	ctx := context.Background()
	// butketName := ""
	filePath := "./static/images.png"
	key := filepath.Base(filePath) + time.Now().Format("20060102150405")
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	client := NewS3Client(ctx)
	client.PutObject(ctx, key, file)
	//	createBucket(ctx, "my-bucket")
	//	getListBucket(ctx)
}

func NewS3Client(ctx context.Context) *S3Client {
	//if s3Client != nil {
	//	return s3Client
	//}
	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	// create an aws s3 client
	s3Client := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	return &S3Client{
		Client: s3Client,
		Bucket: bucket,
		Region: region,
	}
}

// PutObject TODO: put object to s3 and return string path if error not exits
func (s *S3Client) PutObject(ctx context.Context, key string, data *os.File) (string, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   data,
	}
	_, err := s.Client.PutObject(ctx, input)
	if err != nil {
		log.Fatalf("failed to put object, %v", err)
	}
	fmt.Printf("Successfully uploaded file to %s\n", s.Bucket)
	return key, nil
}

//func getListBucket(ctx context.Context) {
//	// get s3 client
//	s3Client := getS3Client(ctx)
//	result, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
//	if err != nil {
//		var ae smithy.APIError
//		if errors.As(err, &ae) && ae.ErrorCode() == "AccessDenied" {
//			fmt.Println("You don't have permission to list buckets for this account.")
//		} else {
//			fmt.Printf("Couldn't list buckets for your account. Here's why: %v\n", err)
//		}
//		return
//	}
//	if len(result.Buckets) == 0 {
//		fmt.Println("You don't have any S3 buckets")
//		return
//	} else {
//		fmt.Println("You have the following S3 buckets:")
//		for _, bucket := range result.Buckets {
//			fmt.Println(*bucket.Name)
//		}
//	}
//
//}

//func createBucket(ctx context.Context, bucketName string) {
//	// get s3 client
//	s3Client := getS3Client(ctx)
//	_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
//		Bucket: &bucketName,
//		CreateBucketConfiguration: &types.CreateBucketConfiguration{
//			LocationConstraint: types.BucketLocationConstraintUsEast2,
//		},
//	})
//	if err != nil {
//		var ae smithy.APIError
//		if errors.As(err, &ae) && ae.ErrorCode() == "BucketAlreadyExists" {
//			fmt.Printf("Bucket %s already exists.\n", bucketName)
//		} else {
//			fmt.Printf("Couldn't create bucket %s. Here's why: %v\n", bucketName, err)
//		}
//		return
//	}
//	fmt.Printf("Successfully created bucket %s\n", bucketName)
//}
