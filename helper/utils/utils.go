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

package utils

import (
	"eks-spot-termination-monitor/helper/logging"
	"os"

	"github.com/joho/godotenv"
)

var logger = logging.Log()

func LoadEnv() {
	err := godotenv.Load("./config/.env")
	Must("Error loading env file", err)
}

func ReadConfig(key string) (string, error) {
	val := os.Getenv(key)
	return val, nil
}

func Must(msg string, err error) {
	if err != nil {
		logger.Fatalf("%s: %v", msg, err)
	}
}

func LogError(msg string, err error) {
	if err != nil {
		logger.Errorf("%s: %v", msg, err)
	}
}
