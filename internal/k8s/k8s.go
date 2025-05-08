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

package k8s

import (
	"context"
	"eks-spot-termination-monitor/helper/logging"
	"eks-spot-termination-monitor/helper/utils"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	logger       = logging.Log()
	k8sClientSet *kubernetes.Clientset
	once         sync.Once
)

type ImpactedPod struct {
	PodName      string
	PodNamespace string
	PodState     string
}

// reference:
// 1. https://github.com/kubernetes/client-go?tab=readme-ov-file#how-to-get-it
// 2. https://github.com/kubernetes/client-go/blob/master/examples/in-cluster-client-configuration/main.go
// 3. https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/main.go

// GetK8sClient Reusable function to get k8s clientset.
/*
1st Call to GetK8sClient()	once.Do() runs the func() {...}, creates config, builds client, stores in k8sClientSet.
2nd Call to GetK8sClient()	once.Do() does nothing. k8sClientSet is already ready, returns it.
3rd Call to GetK8sClient() Same, just returns k8sClientSet.
*/
func GetK8sClient() (*kubernetes.Clientset, error) {
	var err error
	//once.Do() guarantees that the function inside it executes exactly one time,
	//no matter how many times GetK8sClient() is called — even across multiple threads/goroutines.
	once.Do(func() {
		var config *rest.Config
		config, err = rest.InClusterConfig()
		if err != nil {
			kubeconfig, _ := utils.ReadConfig("KUBECONFIG")
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err != nil {
				return
			}
		}
		// save the value to global variable at first time
		k8sClientSet, err = kubernetes.NewForConfig(config)
	})
	// always return value of k8sClientSet from global variable
	return k8sClientSet, err
}

// TrackImpactedPods watches for pod events on a specific node and prints the events to the console
// reference:
// 1. https://github.com/kubernetes/client-go/issues/410#issuecomment-389680017
// 2.https://github.com/kubernetes/client-go/blob/17ad09eb27d15d66be61e55729f6da1ef5986faf/examples/out-of-cluster-client-configuration/main.go#L61
func TrackImpactedPods(nodeName string) (*[]ImpactedPod, error) {
	clientset, err := GetK8sClient()
	if err != nil {
		logger.Panicln("Error creating Kubernetes client", err)
	}

	// Using context with a timeout to avoid hanging indefinitely, otherwise if API is stuck, app will be blocked unnecessarily.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})

	if err != nil {
		return nil, err
	}

	// pre-allocating slice to avoid repeated allocations.
	impactedPods := make([]ImpactedPod, 0, len(pods.Items))

	for _, pod := range pods.Items {
		impactedPods = append(impactedPods, ImpactedPod{
			PodName:      pod.Name,
			PodNamespace: pod.Namespace,
			PodState:     string(pod.Status.Phase),
		})
	}
	return &impactedPods, nil
}
