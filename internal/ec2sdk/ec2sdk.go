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

package ec2sdk

import (
	"context"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

var (
	ec2Client *ec2.Client
	once      sync.Once
)

type EC2Details struct {
	InstanceID     *string
	InstanceType   string
	InstanceArch   string
	InstanceState  string
	PrivateDNSName *string
}

// GetEC2Client GetK8sClient Reusable function to get EC2 clientset.
/*
1st Call to GetEC2Client()	once.Do() runs the func() {...}, creates config, builds client, stores in k8sClientSet.
2nd Call to GetEC2Client()	once.Do() does nothing. k8sClientSet is already ready, returns it.
3rd Call to GetEC2Client() Same, just returns k8sClientSet.
*/

func GetEC2Client() (*ec2.Client, error) {
	var cfg aws.Config
	var err error
	//once.Do() guarantees that the function inside it executes exactly one time,
	//no matter how many times GetEC2Client() is called — even across multiple threads/goroutines.
	once.Do(func() {
		ctx := context.Background()
		profile := os.Getenv("AWS_PROFILE")

		// reference: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws#Config
		// reference: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws#AnonymousCredentials
		if profile != "" {
			// Running from local machine
			cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
		} else {
			// Running from EKS cluster with service account attached with IAM Role
			cfg, err = config.LoadDefaultConfig(ctx)
		}

		if err != nil {
			return
		}

		// Create EC2 client
		// reference: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2#NewFromConfig
		ec2Client = ec2.NewFromConfig(cfg)
	})
	return ec2Client, err
}

// DescribeInstance retrieves information about an EC2 instance using its instance ID.
// reference: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/ec2#Client.DescribeInstances
func DescribeInstance(instanceID string) (*EC2Details, error) {
	client, _ := GetEC2Client()

	input := ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID}}

	response, err := client.DescribeInstances(context.TODO(), &input)
	if err != nil {
		return nil, err
	}

	var ec2Data EC2Details
	for _, reservation := range response.Reservations {
		for _, instance := range reservation.Instances {
			ec2Data = EC2Details{
				InstanceID:     instance.InstanceId,
				InstanceType:   string(instance.InstanceType),
				InstanceArch:   string(instance.Architecture),
				InstanceState:  string(instance.State.Name),
				PrivateDNSName: instance.PrivateDnsName,
			}
		}
	}
	return &ec2Data, nil
}
