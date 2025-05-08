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

type RequestData struct {
	URL    string                   `json:"url"`
	Body   EC2SpotInterruptionEvent `json:"body"`
	Method string
}

type EC2SpotInterruptionEvent struct {
	Version    string    `json:"version"`
	ID         string    `json:"id"`
	DetailType string    `json:"detail-type"`
	Source     string    `json:"source"`
	Account    string    `json:"account"`
	Time       string    `json:"time"`
	Region     string    `json:"region"`
	Resources  []string  `json:"resources"`
	Detail     EC2Detail `json:"detail"`
}

type EC2Detail struct {
	InstanceID     string `json:"instance-id"`
	InstanceAction string `json:"instance-action"`
}

// For future use
//type ProcessedData struct {
//	InstanceID    string `json:"instance-id"`
//	InstanceType  string `json:"instance-type"`
//	InstanceArch  string `json:"instance-arch"`
//	InstanceState string `json:"instance-state"`
//	EC2SPOTEvent  string `json:"ec2sdk-spot-event"`
//	NodeName      string `json:"node-name"`
//	ImpactedPod   string `json:"impacted-pods"`
//	PodNamespace  string `json:"pod-namespace"`
//	PodState      string `json:"pod-state"`
//}
