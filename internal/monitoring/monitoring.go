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

package monitoring

import (
	"eks-spot-termination-monitor/helper/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	logger = logging.Log()

	InterruptionsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eks_spot_interruption_last_seen",
			Help: "Gauge set to 1 when an interruption is seen; can be used with timestamp.",
		},
		[]string{},
	)

	InterruptionsByInstanceTypeOverLifetime = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "eks_spot_interruptions_by_instance_type_over_lifetime",
			Help: "Percentage of spot interruptions per instance type",
		},
		[]string{"instance_type"},
	)

	InterruptionsByInstanceTypeAtGivenTime = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "eks_spot_interruptions_by_instance_type_total",
			Help: "Total number of spot instance interruptions by instance type.",
		},
		[]string{"instance_type"},
	)

	PodsImpacted = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eks_spot_interruptions_pods_impacted_per_termination",
			Help: "Number of pods impacted due to termination",
		},
		[]string{},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(InterruptionsGauge)
	prometheus.MustRegister(InterruptionsByInstanceTypeOverLifetime)
	prometheus.MustRegister(InterruptionsByInstanceTypeAtGivenTime)
	prometheus.MustRegister(PodsImpacted)
}

func ExposeMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":2112", nil)
		if err != nil {
			logger.Panic("Error starting metrics server: ", err)
		}
	}()
}
