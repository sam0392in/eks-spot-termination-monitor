/*
	Copyright 2025 Samarth Kanungo

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package queue

import (
	"context"
	"eks-spot-termination-monitor/helper/logging"
	"eks-spot-termination-monitor/helper/utils"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var logger = logging.Log()

// reference: https://docs.aws.amazon.com/code-library/latest/ug/go_2_sqs_code_examples.html
func getSQSClient() (*sqs.Client, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-west-1"))
	if err != nil {
		logger.Error("Error loading default config, Error: " + err.Error())
		fmt.Println(err)
		return nil, err
	}
	sqsClient := sqs.NewFromConfig(sdkConfig)
	return sqsClient, nil
}

// ReceiveFromSQS GetMessages Receives message from SQS
// Reference: https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/go/sqs/ReceiveMessage/ReceiveMessage.go
func ReceiveFromSQS(queueURL string) (*sqs.ReceiveMessageOutput, error) {
	// Create an SQS service client
	client, err := getSQSClient()

	maxMessageCountStr, _ := utils.ReadConfig("SQS_MAX_MESSAGE_COUNT")
	maxMessageCount, _ := strconv.ParseInt(maxMessageCountStr, 10, 32)

	pollingIntervalStr, _ := utils.ReadConfig("SQS_POLLING_INTERVAL")
	pollingInterval, _ := strconv.ParseInt(pollingIntervalStr, 10, 32)

	visibilityTimeoutStr, _ := utils.ReadConfig("SQS_VISIBILITY_TIMEOUT")
	visibilityTimeout, _ := strconv.ParseInt(visibilityTimeoutStr, 10, 32)

	params := &sqs.ReceiveMessageInput{
		QueueUrl:            &queueURL,
		MaxNumberOfMessages: int32(maxMessageCount),   // sequential processing needed.
		WaitTimeSeconds:     int32(pollingInterval),   // Long polling SQS to reduce number of calls. This will keep calls within the AWS limit.
		VisibilityTimeout:   int32(visibilityTimeout), // this will prevent ReceiptHandle of message from getting expired.
	}
	msgResult, err := client.ReceiveMessage(context.TODO(), params)
	utils.LogError("Error receiving message from SQS", err)
	return msgResult, nil
}

func DeleteMessageFromSQS(queueURL string, messageHandle *string) error {
	// Create an SQS service client
	client, err := getSQSClient()

	params := &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: messageHandle,
	}

	_, err = client.DeleteMessage(context.TODO(), params)
	if err != nil {
		return err
	}
	logger.Info("Message deleted successfully from the queue.")
	return nil
}
