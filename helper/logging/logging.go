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

package logging

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type panicFormatter struct{}

func (f *panicFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[PANIC] %s: %s\n", entry.Message, entry.Error)), nil
}

func Log() *logrus.Logger {
	var Logger *logrus.Logger
	Logger = logrus.New()
	Logger.SetFormatter(&panicFormatter{}) // Use custom formatter for Panic level
	Logger.SetReportCaller(false)          // Disable caller information
	Logger.Formatter = &logrus.JSONFormatter{}
	val := os.Getenv("LOG_LEVEL")
	switch val {
	case "info":
		Logger.Level = logrus.InfoLevel
	case "debug":
		Logger.Level = logrus.DebugLevel
	}
	return Logger
}
