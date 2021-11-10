// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config // import "go.opentelemetry.io/collector/config"

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

// Service defines the configurable components of the service.
type Service struct {
	// Telemetry is the configuration for collector's own telemetry.
	Telemetry ServiceTelemetry `mapstructure:"telemetry"`

	// Extensions are the ordered list of extensions configured for the service.
	Extensions []ComponentID `mapstructure:"extensions"`

	// Pipelines are the set of data pipelines configured for the service.
	Pipelines Pipelines `mapstructure:"pipelines"`
}

// ServiceTelemetry defines the configurable settings for service telemetry.
type ServiceTelemetry struct {
	Logs ServiceTelemetryLogs `mapstructure:"logs"`
}

func (srvT *ServiceTelemetry) validate() error {
	return srvT.Logs.validate()
}

// ServiceTelemetryLogs defines the configurable settings for service telemetry logs.
// This MUST be compatible with zap.Config. Cannot use directly zap.Config because
// the collector uses mapstructure and not yaml tags.
type ServiceTelemetryLogs struct {
	// Level is the minimum enabled logging level.
	Level zapcore.Level `mapstructure:"level"`

	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally.
	Development bool `mapstructure:"development"`

	// Encoding sets the logger's encoding.
	// Valid values are "json" and "console".
	Encoding string `mapstructure:"encoding"`

	// DisableCaller stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`

	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktraces are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace"`

	// Sampling sets a sampling policy. A nil SamplingConfig disables sampling.
	Sampling *SamplingConfig `json:"sampling" yaml:"sampling"`

	// EncoderConfig sets options for the chosen encoder. See
	// zapcore.EncoderConfig for details.
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`

	// OutputPaths is a list of URLs or file paths to write logging output to.
	// See Open for details.
	OutputPaths []string `json:"outputPaths" yaml:"outputPaths"`

	// ErrorOutputPaths is a list of URLs to write internal logger errors to.
	// The default is standard error.
	//
	// Note that this setting only affects internal errors; for sample code that
	// sends error-level logs to a different location from info- and debug-level
	// logs, see the package-level AdvancedConfiguration example.
	ErrorOutputPaths []string `json:"errorOutputPaths" yaml:"errorOutputPaths"`

	// InitialFields is a collection of fields to add to the root logger.
	InitialFields map[string]interface{} `json:"initialFields" yaml:"initialFields"`

}

// Copied from zap/config.go
// SamplingConfig sets a sampling strategy for the logger. Sampling caps the
// global CPU and I/O load that logging puts on your process while attempting
// to preserve a representative subset of your logs.
//
// If specified, the Sampler will invoke the Hook after each decision.
//
// Values configured here are per-second. See zapcore.NewSamplerWithOptions for
// details.
type SamplingConfig struct {
	Initial    int                                           `json:"initial" yaml:"initial"`
	Thereafter int                                           `json:"thereafter" yaml:"thereafter"`
	Hook       func(zapcore.Entry, zapcore.SamplingDecision) `json:"-" yaml:"-"`
}

type EncoderConfig struct {
	// Set the keys used for each log entry. If any key is empty, that portion
	// of the entry is omitted.
	MessageKey    string `json:"messageKey" yaml:"messageKey"`
	LevelKey      string `json:"levelKey" yaml:"levelKey"`
	TimeKey       string `json:"timeKey" yaml:"timeKey"`
	NameKey       string `json:"nameKey" yaml:"nameKey"`
	CallerKey     string `json:"callerKey" yaml:"callerKey"`
	FunctionKey   string `json:"functionKey" yaml:"functionKey"`
	StacktraceKey string `json:"stacktraceKey" yaml:"stacktraceKey"`
	LineEnding    string `json:"lineEnding" yaml:"lineEnding"`
	// Configure the primitive representations of common complex types. For
	// example, some users may want all time.Times serialized as floating-point
	// seconds since epoch, while others may prefer ISO8601 strings.
	EncodeLevel    zapcore.LevelEncoder    `json:"levelEncoder" yaml:"levelEncoder"`
	EncodeTime     zapcore.TimeEncoder     `json:"timeEncoder" yaml:"timeEncoder"`
	EncodeDuration zapcore.DurationEncoder `json:"durationEncoder" yaml:"durationEncoder"`
	EncodeCaller   zapcore.CallerEncoder   `json:"callerEncoder" yaml:"callerEncoder"`
	// Unlike the other primitive type encoders, EncodeName is optional. The
	// zero value falls back to FullNameEncoder.
	EncodeName zapcore.NameEncoder `json:"nameEncoder" yaml:"nameEncoder"`
	// Configures the field separator used by the console encoder. Defaults
	// to tab.
	ConsoleSeparator string `json:"consoleSeparator" yaml:"consoleSeparator"`
}

func (srvTL *ServiceTelemetryLogs) validate() error {
	if srvTL.Encoding != "json" && srvTL.Encoding != "console" {
		return fmt.Errorf(`service telemetry logs invalid encoding: %q, valid values are "json" and "console"`, srvTL.Encoding)
	}
	return nil
}

// DataType is a special Type that represents the data types supported by the collector. We currently support
// collecting metrics, traces and logs, this can expand in the future.
type DataType = Type

// Currently supported data types. Add new data types here when new types are supported in the future.
const (
	// TracesDataType is the data type tag for traces.
	TracesDataType DataType = "traces"

	// MetricsDataType is the data type tag for metrics.
	MetricsDataType DataType = "metrics"

	// LogsDataType is the data type tag for logs.
	LogsDataType DataType = "logs"
)

// Pipeline defines a single pipeline.
type Pipeline struct {
	Receivers  []ComponentID `mapstructure:"receivers"`
	Processors []ComponentID `mapstructure:"processors"`
	Exporters  []ComponentID `mapstructure:"exporters"`
}

// Pipelines is a map of names to Pipelines.
type Pipelines map[ComponentID]*Pipeline
