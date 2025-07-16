#!/bin/bash

# Wait for LocalStack to be ready
# echo "Waiting for LocalStack to be ready..."
# while ! curl -s http://localhost:4566/_localstack/health | grep -q '"sqs": "running"'; do
#     sleep 1
# done

# Create SQS queue
echo "Creating SQS queue..."
aws --endpoint-url=http://localhost:4566 sqs create-queue \
    --queue-name audit-log-queue \
    --attributes '{
        "VisibilityTimeout": "30",
        "MessageRetentionPeriod": "86400",
        "DelaySeconds": "0",
        "ReceiveMessageWaitTimeSeconds": "20"
    }'

echo "LocalStack initialization complete!" 