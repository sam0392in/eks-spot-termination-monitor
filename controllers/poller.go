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

package controllers

import (
	"eks-spot-termination-monitor/helper/logging"
	"eks-spot-termination-monitor/helper/utils"
	"eks-spot-termination-monitor/internal/monitoring"
	"eks-spot-termination-monitor/internal/queue"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

var (
	logger      = logging.Log()
	sqsMessages = make(chan types.Message)
	wg          sync.WaitGroup
)

// Wait blocks until all wg.Done() are called
func Wait() {
	wg.Wait()
}

func PollQueue(stopChan chan struct{}) {
	queueURL, err := utils.ReadConfig("QUEUE_URL")
	logger.Info("Starting PollQueue...")
	utils.LogError("Error reading queue URL from config", err)
	clientPollFrequencyStr, _ := utils.ReadConfig("CLIENT_POLL_FREQUENCY")
	clientPollFrequency, _ := strconv.ParseInt(clientPollFrequencyStr, 10, 64)

	for {
		select {
		case <-stopChan:
			logger.Info("Stopping PollQueue...")
			close(sqsMessages) // Close the channel to signal the executors to stop
			return
		default:
			msg, err := queue.ReceiveFromSQS(queueURL)
			utils.LogError("Error receiving message from SQS", err)

			// do not process empty messages
			if len(msg.Messages) == 0 {
				time.Sleep(time.Duration(clientPollFrequency) * time.Second) // Sleep if no messages
				logger.Debug("No messages in queue, next poll in " + clientPollFrequencyStr + " seconds")

				// Reset metrics when no event is received
				monitoring.InterruptionsGauge.WithLabelValues().Set(0)
				monitoring.PodsImpacted.WithLabelValues().Set(0)
				monitoring.InterruptionsByInstanceTypeAtGivenTime.Reset()
				continue
			}

			for _, message := range msg.Messages {
				// Send each message separately to the channel
				logger.Debug("Message sent to channel")
				wg.Add(1)
				sqsMessages <- message
			}
		}
	}
}
