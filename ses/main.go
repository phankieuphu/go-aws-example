package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

const (
	NumberWorker = 3
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	Sender = "phankieuphu@example.com"

	// Replace recipient@example.com with a "To" address. If your account
	// is still in the sandbox, this address must be verified.
	Recipient = "recipient@example.com"

	// Specify a configuration set. To use a configuration
	// set, comment the next line and line 92.
	//ConfigurationSet = "ConfigSet"

	// The subject line for the email.
	Subject = "Amazon SES Test (AWS SDK for Go)"

	// The HTML body for the email.
	HtmlBody = "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
		"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
		"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"

	//The email body for recipients with non-HTML email clients.
	TextBody = "This email was sent with Amazon SES using the AWS SDK for Go."

	// The character encoding for the email.
	CharSet = "UTF-8"
)

func main() {
	// Load AWS configuration with explicit credentials
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				// Use LocalStack SES endpoint
				if service == ses.ServiceID {
					return aws.Endpoint{
						URL:           "http://localhost:4566",
						SigningRegion: "us-east-1",
					}, nil
				}
				return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
			},
		)),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("test", "test", ""), // Dummy credentials
		),
	)
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	// create ses client
	svc := ses.NewFromConfig(cfg)

	// init context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	lenEmail := 10

	emailChannel := make(chan int, lenEmail)
	var wg sync.WaitGroup

	for i := 1; i <= NumberWorker; i++ {
		wg.Add(1)
		go worker(ctx, emailChannel, svc, &wg)
	}

	for email := 1; email <= lenEmail; email++ {
		emailChannel <- email
	}
	close(emailChannel)
	wg.Wait()
	log.Println("All emails sent")

}

func worker(ctx context.Context, emails <-chan int, svc *ses.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	// handle goroutines and channel here
	for email := range emails {
		sendEmail(svc, email, ctx)
	}
}

func sendEmail(svc *ses.Client, index int, ctx context.Context) {
	select {
	case <-ctx.Done():
		log.Fatalf("context cancelled: %v", ctx.Err())
	default:
		sender := "phankieuphu@gmail.com"
		recipient := "recipient@gmail.com"

		input := ses.SendEmailInput{
			Destination: &types.Destination{ToAddresses: []string{recipient}},
			// Destination: &ses.Destination{ToAddresses: []*string{aws.String(recipient)}},
			Message: &types.Message{
				Subject: &types.Content{
					Data:    aws.String(string(index) + " " + Subject),
					Charset: aws.String(CharSet),
				},
				Body: &types.Body{
					Html: &types.Content{Data: aws.String(HtmlBody), Charset: aws.String(CharSet)},
					Text: &types.Content{Data: aws.String(TextBody), Charset: aws.String(CharSet)},
				},
			},
			Source: aws.String(sender),
		}

		_, err := svc.SendEmail(ctx, &input)
		if err != nil {
			log.Fatalf("failed to send email: %v", err)
		}
		fmt.Println("Email Sent to address: " + recipient)
		// handle channel here
	}
	// handle sending email here

}
