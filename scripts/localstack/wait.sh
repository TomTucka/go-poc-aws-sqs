#! /usr/bin/env bash

echo "Waiting for SQS queues"

iterations=0

while [ "$iterations" -lt 300 ]
do
  queues=$(awslocal sqs list-queues)

  if [[ $queues = *'"http://localstack:4566/000000000000/reindex.fifo"'* ]]
  then
    echo "Found all expected SQS queues after $iterations seconds"
    exit 0
  fi

  ((iterations++))
  sleep 1
done

echo "Waited $iterations seconds for SQS queues before giving up"
echo "sqs list-queues results:"
echo "----------------------------------"
echo "$queues"
echo "----------------------------------"

exit 1