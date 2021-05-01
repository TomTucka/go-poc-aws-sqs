#! /usr/bin/env bash

awslocal secretsmanager create-secret --name local/test-secret \
   --description "Local Test Secret" \
   --secret-string "ThisIsATestSecret"

awslocal sqs create-queue --queue-name reindex.fifo --attributes FifoQueue=true,ContentBasedDeduplication=true,VisibilityTimeout=30,ReceiveMessageWaitTimeSeconds=0

for iterations in {1..4000}
do
   awslocal sqs send-message --queue-url http://localstack:4566/000000000000/reindex.fifo --message-body "file:///scripts/test.json" --message-group-id "$iterations"
done


# awslocal sqs get-queue-attributes --queue-url http://localstack:4566/000000000000/reindex.fifo --attribute-names All