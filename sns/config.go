package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func GetSNSClient(ctx context.Context) *sns.Client {
	if snsInstance != nil {
		return snsInstance
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		fmt.Println("unable to load SDK config, " + err.Error())
		return nil
	}

	snsClient := sns.NewFromConfig(cfg, func(o *sns.Options) {
		o.EndpointResolver = sns.EndpointResolverFromURL("http://localhost:4566")
	})
	return snsClient
}
