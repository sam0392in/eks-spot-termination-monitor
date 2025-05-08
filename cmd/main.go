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
	"eks-spot-termination-monitor/controllers"
	"eks-spot-termination-monitor/helper/logging"
	"eks-spot-termination-monitor/helper/utils"
	"eks-spot-termination-monitor/internal/monitoring"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	logger                            = logging.Log()
	shutDownGracePeriod time.Duration = 10 * time.Second
)

func main() {
	utils.LoadEnv()
	monitoring.RegisterMetrics()
	monitoring.ExposeMetrics()

	stopChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	workerCount, _ := utils.ReadConfig("GO_WORKERS")
	workerPool, _ := strconv.Atoi(workerCount)

	// Start PollQueue
	go controllers.PollQueue(stopChan)

	// Start Executor Pool
	for i := 1; i <= workerPool; i++ {
		go controllers.Executor(i)
	}

	<-signalChan

	logger.Info("Received shutdown signal")

	close(stopChan)
	controllers.Wait()

	// Give the controller some time to shut down gracefull
	time.Sleep(shutDownGracePeriod)
	logger.Info("Controller stopped")
}
