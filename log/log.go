// Copyright 2019 Adverax. All Rights Reserved.
// This file is part of project
//
//      http://github.com/adverax/echo
//
// Licensed under the MIT (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://github.com/adverax/echo/blob/master/LICENSE
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"fmt"
	"io"
	"io/ioutil"
	syslog "log"
	"os"
	"strings"
	"sync/atomic"
)

const (
	ClassTrace Class = iota + 1
	ClassInfo
	ClassWarning
	ClassError
)

const (
	delimiter = "\n"
)

type Class uint8

// Statistics for logger
type Metrics struct {
	Traces   int32 `json:"traces,omitempty"`   // Count of trace messages
	Infos    int32 `json:"infos,omitempty"`    // Count of information messages
	Errors   int32 `json:"errors,omitempty"`   // Count of error messages
	Warnings int32 `json:"warnings,omitempty"` // Count of warning messages
}

type Logger interface {
	Trace(v interface{})
	Info(v interface{})
	Warning(v interface{})
	Error(v interface{})
	Metrics() Metrics
}

type logger struct {
	level   int
	prefix  string // Common prefix for all types
	trace   *syslog.Logger
	info    *syslog.Logger
	warning *syslog.Logger
	error   *syslog.Logger
	metrics Metrics
}

func (log *logger) Trace(v interface{}) {
	if s, ok := valueToString(v); ok {
		atomic.AddInt32(&log.metrics.Traces, 1)
		_ = log.trace.Output(2, delimiter+s)
	}
}

func (log *logger) Info(v interface{}) {
	if s, ok := valueToString(v); ok {
		atomic.AddInt32(&log.metrics.Infos, 1)
		_ = log.info.Output(2, delimiter+s)
	}
}

func (log *logger) Warning(v interface{}) {
	if s, ok := valueToString(v); ok {
		atomic.AddInt32(&log.metrics.Warnings, 1)
		_ = log.warning.Output(2, delimiter+s)
	}
}

func (log *logger) Error(v interface{}) {
	if s, ok := valueToString(v); ok {
		atomic.AddInt32(&log.metrics.Errors, 1)
		_ = log.error.Output(2, delimiter+s)
	}
}

func (log *logger) Metrics() Metrics {
	return Metrics{
		Traces:   atomic.LoadInt32(&log.metrics.Traces),
		Infos:    atomic.LoadInt32(&log.metrics.Infos),
		Warnings: atomic.LoadInt32(&log.metrics.Warnings),
		Errors:   atomic.LoadInt32(&log.metrics.Errors),
	}
}

func New(
	trace io.Writer,
	info io.Writer,
	warning io.Writer,
	error io.Writer,
	prefix string, // Prefix for each entry.
) Logger {
	return NewEx(
		trace,
		info,
		warning,
		error,
		syslog.Ldate|syslog.Ltime,
		prefix,
		true,
	)
}

func NewEx(
	trace io.Writer,
	info io.Writer,
	warning io.Writer,
	error io.Writer,
	flag int, // System flags
	prefix string, // Prefix for each entry
	labels bool, // Enable type labels for each entry.
) Logger {
	if labels {
		return &logger{
			trace:   syslog.New(trace, prefix+"TRACE: ", flag),
			info:    syslog.New(info, prefix+"INFO: ", flag),
			warning: syslog.New(warning, prefix+"WARNING: ", flag),
			error:   syslog.New(error, prefix+"ERROR: ", flag),
		}
	}

	return &logger{
		trace:   syslog.New(trace, prefix, flag),
		info:    syslog.New(info, prefix, flag),
		warning: syslog.New(warning, prefix, flag),
		error:   syslog.New(error, prefix, flag),
	}
}

var (
	// Stub without any logging
	discardLogger = New(
		ioutil.Discard,
		ioutil.Discard,
		ioutil.Discard,
		ioutil.Discard,
		"",
	)
)

func NewDiscard() Logger {
	return discardLogger
}

// Create new debug logger with prefix for each entry.
func NewDebug(prefix string) Logger {
	return New(
		os.Stdout,
		os.Stdout,
		os.Stdout,
		os.Stderr,
		prefix,
	)
}

// Create new file logger with prefix for each entry.
func NewFile(
	fileName string,
	prefix string,
) Logger {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf("failed to open log file %q: %v", fileName, err))
	}
	return New(
		ioutil.Discard,
		ioutil.Discard,
		file,
		file,
		prefix,
	)
}

func Append(logger Logger, v interface{}) {
	if logger != nil && v != nil {
		logger.Error(v)
	}
}

func valueToString(v interface{}) (string, bool) {
	switch val := v.(type) {
	case error:
		return fmt.Sprintf("%+v\n", val), true
	case string:
		return val + "\n", true
	default:
		return "", false
	}
}

var decoders = map[Class]string{
	ClassTrace:   "trace",
	ClassInfo:    "info",
	ClassWarning: "warning",
	ClassError:   "error",
}

var encoders = map[string]Class{
	"trace":   ClassTrace,
	"info":    ClassInfo,
	"warning": ClassWarning,
	"error":   ClassError,
}

func DecodeClassName(class Class) string {
	if res, found := decoders[class]; found {
		return res
	}
	return "error"
}

func EncodeClassName(class string) Class {
	if res, found := encoders[strings.ToLower(class)]; found {
		return res
	}
	return ClassError
}
