package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSConfig struct {
	Region          string `mapstructure:"region"`
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	QueueURL        string `mapstructure:"queue_url"`
}

func DefaultSQSConfig() *SQSConfig {
	return &SQSConfig{
		Region:          getEnvOrDefault("AWS_REGION", "us-east-1"),
		Endpoint:        getEnvOrDefault("AWS_SQS_ENDPOINT", "http://localhost:4566"), // LocalStack default
		AccessKeyID:     getEnvOrDefault("AWS_ACCESS_KEY_ID", "dummy"),
		SecretAccessKey: getEnvOrDefault("AWS_SECRET_ACCESS_KEY", "dummy"),
		QueueURL:        getEnvOrDefault("AWS_SQS_QUEUE_URL", "http://localhost:4566/000000000000/audit-log-queue"),
	}
}

func (c *SQSConfig) GetClient() (*sqs.Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == sqs.ServiceID {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           c.Endpoint,
				SigningRegion: c.Region,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(c.Region),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			c.AccessKeyID,
			c.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return sqs.NewFromConfig(cfg), nil
}
