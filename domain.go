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

package echo

import (
	stdContext "context"
	"fmt"
	"net/url"
	"time"

	"github.com/adverax/echo/data"
)

// Abstract data storage
type Storage interface {
	// Get value by key.
	Get(key string, dst interface{}) error
	// Set value with key and expire time.
	Set(key string, val interface{}, timeout time.Duration) error
	// Check if value exists or not.
	IsExist(key string) (bool, error)
	// Delete cached value by key.
	Delete(key string) error
}

// Locale represents localization strategy.
type Locale interface {
	// Format date at the current location
	FormatDate(t time.Time) string
	// Format time at the current location
	FormatTime(t time.Time) string
	// Format datetime at the current location
	FormatDateTime(t time.Time) string

	// Parse date at the current location
	ParseDate(value string) (time.Time, error)
	// Parse time at the current location
	ParseTime(value string) (time.Time, error)
	// Parse datetime at the current location
	ParseDateTime(value string) (time.Time, error)

	// Get active language identifier
	Language() uint16
	// Get active timezone identifier
	Timezone() uint16
	// Get active location
	Location() *time.Location

	// Get message translation into the current language
	Message(ctx stdContext.Context, id uint32) (string, error)
	// Get resource translation into the current language
	Resource(ctx stdContext.Context, id uint32) (string, error)
	// Get data source translation into the current language
	DataSet(ctx stdContext.Context, id uint32) (data.Set, error)
}

// MessageFamily for selected language
type MessageFamily interface {
	Fetch(ctx stdContext.Context, id uint32) (string, error)
}

// ResourceFamily for selected language
type ResourceFamily interface {
	Fetch(ctx stdContext.Context, id uint32) (string, error)
}

// DataSetFamily for selected language
type DataSetFamily interface {
	Fetch(ctx stdContext.Context, id uint32) (data.Set, error)
}

// BaseLocale is a simple Locale structure.
type BaseLocale struct {
	DateFormat     string
	TimeFormat     string
	DateTimeFormat string

	Lang  uint16 // Language identifier
	TZone uint16 // Timezone identifier
	Loc   *time.Location

	Messages  MessageFamily
	Resources ResourceFamily
	DataSets  DataSetFamily
}

func (loc *BaseLocale) Language() uint16 {
	return loc.Lang
}

func (loc *BaseLocale) Timezone() uint16 {
	return loc.TZone
}

func (loc *BaseLocale) Location() *time.Location {
	return loc.Loc
}

func (loc *BaseLocale) FormatDate(t time.Time) string {
	return t.In(loc.Loc).Format(loc.DateFormat)
}

func (loc *BaseLocale) FormatTime(t time.Time) string {
	return t.In(loc.Loc).Format(loc.TimeFormat)
}

func (loc *BaseLocale) FormatDateTime(t time.Time) string {
	return t.In(loc.Loc).Format(loc.DateTimeFormat)
}

func (loc *BaseLocale) ParseDate(value string) (time.Time, error) {
	return time.ParseInLocation(loc.DateFormat, value, loc.Loc)
}

func (loc *BaseLocale) ParseTime(value string) (time.Time, error) {
	return time.ParseInLocation(loc.TimeFormat, value, loc.Loc)
}

func (loc *BaseLocale) ParseDateTime(value string) (time.Time, error) {
	return time.ParseInLocation(loc.DateTimeFormat, value, loc.Loc)
}

func (loc *BaseLocale) Message(ctx stdContext.Context, id uint32) (string, error) {
	return loc.Messages.Fetch(ctx, id)
}

func (loc *BaseLocale) Resource(ctx stdContext.Context, id uint32) (string, error) {
	return loc.Resources.Fetch(ctx, id)
}

func (loc *BaseLocale) DataSet(ctx stdContext.Context, id uint32) (data.Set, error) {
	return loc.DataSets.Fetch(ctx, id)
}

// Validation error can be translated into target language.
type ValidationError interface {
	error
	Translate(ctx Context) (string, error)
}

type simpleValidationError struct {
	id uint32
}

func (e *simpleValidationError) Error() string {
	return "Validation error"
}

func (e *simpleValidationError) Translate(
	ctx Context,
) (string, error) {
	return ctx.Locale().Message(ctx, e.id)
}

func NewValidationError(id uint32) ValidationError {
	return &simpleValidationError{id: id}
}

// Complex validation error
type Cause struct {
	Msg  uint32        // Identifier of message
	Text string        // Default literal representation
	Args []interface{} // Custom arguments
}

func (cause *Cause) Error() string {
	if cause.Text != "" {
		if len(cause.Args) == 0 {
			return cause.Text
		}
		return fmt.Sprintf(cause.Text, cause.Args...)
	}

	return fmt.Sprintf("Error %d", uint32(cause.Msg))
}

func (cause *Cause) Translate(
	ctx Context,
) (string, error) {
	msg, err := ctx.Locale().Message(ctx, cause.Msg)
	if err != nil {
		return "", err
	}
	if msg == "" {
		msg = cause.Text
	}

	if len(cause.Args) == 0 {
		return msg, nil
	}

	return fmt.Sprintf(msg, cause.Args...), nil
}

func NewValidationErrorMustBeNotBelow(
	limit string,
) ValidationError {
	return &Cause{
		Msg:  MessageMustBeNotBelow,
		Text: "Value must be not below than %s",
		Args: []interface{}{limit},
	}
}

func NewValidationErrorMustBeNotAbove(
	limit string,
) ValidationError {
	return &Cause{
		Msg:  MessageMustBeNotAbove,
		Text: "Value must be not above than %s",
		Args: []interface{}{limit},
	}
}

func AppendValidationError(
	errs ValidationErrors,
	err error,
) (ValidationErrors, error) {
	if e, ok := err.(ValidationErrors); ok {
		return append(errs, e...), nil
	}

	return errs, err
}

// List of validation errors
// Can be used as error
type ValidationErrors []ValidationError

func (errs ValidationErrors) Error() string {
	return "Validation errors"
}

// Check for empty value
type Empty interface {
	IsEmpty(value interface{}) bool
}

// Url manager (linker)
type UrlLinker interface {
	// Render url
	Render(ctx Context, url *url.URL) (string, error)
	// Expand url by current shard
	Expand(ctx Context, url string) string
	// Collapse url by removing current shard
	Collapse(ctx Context, url string) string
}

var (
	MessageInvalidValue   uint32 = 1
	MessageRequiredValue  uint32 = 2
	MessageMustBeNotBelow uint32 = 3
	MessageMustBeNotAbove uint32 = 4

	ValidationErrorInvalidValue  = NewValidationError(MessageInvalidValue)
	ValidationErrorRequiredValue = NewValidationError(MessageRequiredValue)
)
