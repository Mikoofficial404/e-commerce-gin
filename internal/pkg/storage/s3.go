package storage

import (
	"context"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ConnectS3() *s3.Client {
	ctx := context.Background()
	creditAWS := credentials.NewStaticCredentialsProvider("admin", "passwordRahasia", "")

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(creditAWS),
	)

	if err != nil {
		logrus.Fatal("Gagal membuat config AWS: ", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("http://localhost:9000")
		o.UsePathStyle = true
	})
	return client
}
