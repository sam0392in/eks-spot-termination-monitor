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

package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func main() {

	//------- TEST k8s.go -------//
	//nodeName := "ip-172-23-12-229.eu-west-1.compute.internal"
	//Pods, err := k8s.TrackImpactedPods(nodeName)
	//utils.LogError("Error fetching pods:", err)
	//
	//for _, pod := range *Pods {
	//	println("Pod Name: ", pod.PodName)
	//	println("Pod Namespace: ", pod.PodNamespace)
	//}
	//---------------------------//

	//------- TEST ec2sdk.go -------//
	//instanceID := "i-0a8d3d62a29f67588"
	//data, err := ec2sdk.DescribeInstance(instanceID)
	//utils.LogError("Error fetching instance data:", err)
	//
	//fmt.Println("Instance ID: ", *data.InstanceID)
	//fmt.Println("Instance Type: ", data.InstanceType)
	//fmt.Println("Node Name: ", *data.PrivateDNSName)
	//fmt.Println("Instance Arch: ", data.InstanceArch)
	//fmt.Println("Instance State: ", data.InstanceState)

	//---------------------------//

	// ------- TEST executor.go ------- //
	sampleBody := `{
	  "version": "0",
	  "id": "abcd1234-abcd-1234-abcd-1234abcd1234",
	  "detail-type": "EC2 Spot Instance Interruption Warning",
	  "source": "aws.ec2",
	  "account": "123456789123",
	  "time": "2025-04-28T12:34:56Z",
	  "region": "eu-west-1",
	  "resources": ["arn:aws:ec2:eu-west-1:123456789123:instance/i-0ee6c0327daf387ce"],
	  "detail": {
		"instance-id": "i-0ee6c0327daf387ce",
		"instance-action": "terminate"
	  }
	}`

	msg := types.Message{
		Body: &sampleBody,
	}
	//Pass this in SQS QUeue:
	fmt.Println("Message body: ", msg.Body)
	//controllers.ProcessMessage(msg)
	//----------------------------//

}
