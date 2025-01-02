package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

var (
	S3Client *s3.Client
	Region   = "us-east-1"
)

func main() {
	// create root context
	ctx := context.Background()
	createBucket(ctx, "my-bucket")
	getListBucket(ctx)
}

func getS3Client(ctx context.Context) *s3.Client {
	if S3Client != nil {
		return S3Client
	}
	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(Region))
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	// create an aws s3 client
	return s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.UsePathStyle = true
		o.EndpointResolver = s3.EndpointResolverFromURL("http://localhost:4566")
	})
}

func getListBucket(ctx context.Context) {
	// get s3 client
	s3Client := getS3Client(ctx)
	result, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) && ae.ErrorCode() == "AccessDenied" {
			fmt.Println("You don't have permission to list buckets for this account.")
		} else {
			fmt.Printf("Couldn't list buckets for your account. Here's why: %v\n", err)
		}
		return
	}
	if len(result.Buckets) == 0 {
		fmt.Println("You don't have any S3 buckets")
		return
	} else {
		fmt.Println("You have the following S3 buckets:")
		for _, bucket := range result.Buckets {
			fmt.Println(*bucket.Name)
		}
	}

}

func createBucket(ctx context.Context, bucketName string) {
	// get s3 client
	s3Client := getS3Client(ctx)
	_, err := s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucketName,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraintUsEast2,
		},
	})
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) && ae.ErrorCode() == "BucketAlreadyExists" {
			fmt.Printf("Bucket %s already exists.\n", bucketName)
		} else {
			fmt.Printf("Couldn't create bucket %s. Here's why: %v\n", bucketName, err)
		}
		return
	}
	fmt.Printf("Successfully created bucket %s\n", bucketName)
}
