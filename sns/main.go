package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SnsActions struct {
	SnsClient *sns.Client
}

var (
	snsInstance *sns.Client
)

func main() {
	ctx := context.Background()

}
