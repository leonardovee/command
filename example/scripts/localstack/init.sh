#!/bin/bash
awslocal sns create-topic --name command-events 
awslocal sqs create-queue --queue-name command-events-queue
awslocal sns subscribe \
    --topic-arn arn:aws:sns:us-east-1:000000000000:command-events \
    --protocol sqs \
    --notification-endpoint arn:aws:sqs:us-east-1:000000000000:command-events-queue
