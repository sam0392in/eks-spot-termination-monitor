package controllers

import (
	"eks-spot-termination-monitor/helper/utils"
	"eks-spot-termination-monitor/internal/ec2sdk"
	"eks-spot-termination-monitor/internal/k8s"
	"eks-spot-termination-monitor/internal/monitoring"
	"eks-spot-termination-monitor/internal/queue"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	logging "github.com/sirupsen/logrus"
	"time"
)

func Executor(id int) {
	var queueURL, _ = utils.ReadConfig("QUEUE_URL")
	logger.Info(fmt.Sprintf("Starting executor #%d", id))
	for message := range sqsMessages {
		if processMessage(message) {
			// Delete the message from the queue
			err := queue.DeleteMessageFromSQS(queueURL, message.ReceiptHandle)
			utils.LogError("Error deleting message from SQS:", err)
		} else {
			logger.Info("Message not processed, skipping deletion.")
		}
		wg.Done()
	}
	logger.Info(fmt.Sprintf("Executor #%d shutting down (channel closed)", id))
}

func processMessage(msg types.Message) bool {
	// Process the message
	var result EC2SpotInterruptionEvent

	logger.Debug("Message body:     ", *msg.Body)
	err := json.Unmarshal([]byte(*msg.Body), &result)
	utils.LogError("Error unmarshalling JSON:", err)

	// Fetch EC2 instance data
	ec2Details, err := ec2sdk.DescribeInstance(result.Detail.InstanceID)
	if err != nil {
		logger.Error("Error fetching instance data: ", err)
		return false
	}
	if ec2Details.PrivateDNSName == nil {
		logger.Info(result.Detail.InstanceID + " is already terminated, cannot fetch further details, removing message from queue.")
		return true
	}
	// Track time-based interruption prometheus metrics
	monitoring.InterruptionsGauge.WithLabelValues().Set(1)

	// Instance type-based metrics
	monitoring.InterruptionsByInstanceTypeOverLifetime.WithLabelValues(string(ec2Details.InstanceType)).Inc()
	monitoring.InterruptionsByInstanceTypeAtGivenTime.WithLabelValues(string(ec2Details.InstanceType)).Inc()

	lastEventTime = time.Now()

	// Fetch Impacted Pods
	pods, err := k8s.TrackImpactedPods(*ec2Details.PrivateDNSName)
	if err != nil {
		logger.Error("Error fetching pods:", err)
		return false
	}
	// Track count of impacted pods per termination prometheus metrics
	monitoring.PodsImpacted.WithLabelValues().Set(float64(len(*pods)))

	for _, pod := range *pods {
		logger.WithFields(logging.Fields{
			"InstanceID":    result.Detail.InstanceID,
			"InstanceType":  ec2Details.InstanceType,
			"InstanceArch":  ec2Details.InstanceArch,
			"InstanceState": ec2Details.InstanceState,
			"EC2SPOTEvent":  result.Detail.InstanceAction,
			"NodeName":      *ec2Details.PrivateDNSName,
			"ImpactedPod":   pod.PodName,
			"PodNamespace":  pod.PodNamespace,
			"PodState":      pod.PodState,
		}).Info("SPOT termination captured for " + *ec2Details.PrivateDNSName)
	}

	return true
}
