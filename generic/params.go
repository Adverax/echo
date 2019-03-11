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

package generic

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

var ErrNoMatch = errors.New("no rows")

type GetterFunc func(ctx context.Context, key string, defValue interface{}) (interface{}, error)

func (fn GetterFunc) Get(ctx context.Context, key string, defValue interface{}) (interface{}, error) {
	return fn(ctx, key, defValue)
}

type Getter interface {
	Get(ctx context.Context, key string, defValue interface{}) (interface{}, error)
}

type Setter interface {
	Set(ctx context.Context, key string, value interface{}) error
}

type GetterSetter interface {
	Getter
	Setter
}

type Reader interface {
	Get(ctx context.Context, key string, defValue interface{}) (interface{}, error)

	GetInt(ctx context.Context, key string, defValue int) (int, error)
	GetInt8(ctx context.Context, key string, defValue int8) (int8, error)
	GetInt16(ctx context.Context, key string, defValue int16) (int16, error)
	GetInt32(ctx context.Context, key string, defValue int32) (int32, error)
	GetInt64(ctx context.Context, key string, defValue int64) (int64, error)

	GetUint(ctx context.Context, key string, defValue uint) (uint, error)
	GetUint8(ctx context.Context, key string, defValue uint8) (uint8, error)
	GetUint16(ctx context.Context, key string, defValue uint16) (uint16, error)
	GetUint32(ctx context.Context, key string, defValue uint32) (uint32, error)
	GetUint64(ctx context.Context, key string, defValue uint64) (uint64, error)

	GetFloat32(ctx context.Context, key string, defValue float32) (float32, error)
	GetFloat64(ctx context.Context, key string, defValue float64) (float64, error)

	GetString(ctx context.Context, key string, defValue string) (string, error)
	GetBoolean(ctx context.Context, key string, defValue bool) (bool, error)
	GetTime(ctx context.Context, key string, defValue time.Time) (time.Time, error)
}

type Writer interface {
	Set(ctx context.Context, key string, value interface{}) error

	SetInt(ctx context.Context, key string, value int) error
	SetInt8(ctx context.Context, key string, value int8) error
	SetInt16(ctx context.Context, key string, value int16) error
	SetInt32(ctx context.Context, key string, value int32) error
	SetInt64(ctx context.Context, key string, value int64) error

	SetUint(ctx context.Context, key string, value uint) error
	SetUint8(ctx context.Context, key string, value uint8) error
	SetUint16(ctx context.Context, key string, value uint16) error
	SetUint32(ctx context.Context, key string, value uint32) error
	SetUint64(ctx context.Context, key string, value uint64) error

	SetFloat32(ctx context.Context, key string, value float32) error
	SetFloat64(ctx context.Context, key string, value float64) error

	SetString(ctx context.Context, key string, value string) error
	SetBoolean(ctx context.Context, key string, value bool) error
	SetTime(ctx context.Context, key string, value time.Time) error
}

type ReaderWriter interface {
	Reader
	Writer
}

type Parameters interface {
	ReaderWriter
	JSON() string
}

// List of various named parameters
type Params map[string]interface{}

func (params Params) Get(ctx context.Context, key string, defValue interface{}) (interface{}, error) {
	if value, found := params[key]; found {
		return value, nil
	}
	return defValue, nil
}

func (params Params) GetInt(ctx context.Context, key string, defValue int) (int, error) {
	if value, found := params[key]; found {
		if val, ok := ConvertToInt(value); ok {
			return val, nil
		}
	}
	return defValue, nil
}

func (params Params) GetInt8(ctx context.Context, key string, defValue int8) (int8, error) {
	if value, found := params[key]; found {
		if val, ok := ConvertToInt8(value); ok {
			return val, nil
		}
	}
	return defValue, nil
}

func (params Params) GetInt16(ctx context.Context, key string, defValue int16) (int16, error) {
	if value, found := params[key]; found {
		if val, ok := ConvertToInt16(value); ok {
			return val, nil
		}
	}
	return defValue, nil
}

func (params Params) GetInt32(ctx context.Context, key string, defValue int32) (int32, error) {
	if value, found := params[key]; found {
		if val, ok := ConvertToInt32(value); ok {
			return val, nil
		}
	}
	return defValue, nil
}

func (params Params) GetInt64(ctx context.Context, key string, defValue int64) (int64, error) {
	if value, found := params[key]; found {
		if val, ok := ConvertToInt64(value); ok {
			return val, nil
		}
	}
	return defValue, nil
}

func (params Params) GetUint(ctx context.Context, key string, defValue uint) (uint, error) {
	if value, found := params[key]; found {
		if val, ok := ConvertToUint(value); ok {
			return val, nil
		}
	}
	return defValue, nil
}

func (params Params) GetUint8(ctx context.Context, key string, defValue uint8) (uint8, error) {
	if value, found := params[key]; found {
		if val, ok := ConvertToUint8(value); ok {
			return val, nil
		}
	}
	return defValue, nil
}

func (params Params) GetUint16(ctx context.Context, key string, defValue uint16) (uint16, error) {
	if value, found := params[key]; found {
		if val, ok := ConvertToUint16(value); ok {
			return val, nil
		}
	}
	return defValue, nil
}

func (params Params) GetUint32(ctx context.Context, key string, defValue uint32) (uint32, error) {
	if value, found := params[key]; found {
		if val, ok := ConvertToUint32(value); ok {
			return val, nil
		}
	}
	return defValue, nil
}

func (params Params) GetUint64(ctx context.Context, key string, defValue uint64) (uint64, error) {
	if value, found := params[key]; found {
		if val, ok := ConvertToUint64(value); ok {
			return val, nil
		}
	}
	return defValue, nil
}

func (params Params) GetFloat32(ctx context.Context, key string, defValue float32) (float32, error) {
	if value, found := params[key]; found {
		if val, ok := ConvertToFloat32(value); ok {
			return val, nil
		}
	}
	return defValue, nil
}

func (params Params) GetFloat64(ctx context.Context, key string, defValue float64) (float64, error) {
	if value, found := params[key]; found {
		if val, ok := ConvertToFloat64(value); ok {
			return val, nil
		}
	}
	return defValue, nil
}

func (params Params) GetString(ctx context.Context, key string, defValue string) (string, error) {
	if value, found := params[key]; found {
		if val, ok := ConvertToString(value); ok {
			return val, nil
		}
	}
	return defValue, nil
}

func (params Params) GetBoolean(ctx context.Context, key string, defValue bool) (bool, error) {
	if value, found := params[key]; found {
		if val, ok := ConvertToBoolean(value); ok {
			return val, nil
		}
	}
	return defValue, nil
}

func (params Params) GetTime(ctx context.Context, key string, defValue time.Time) (time.Time, error) {
	if value, found := params[key]; found {
		if val, ok := ConvertToTime(value); ok {
			return val, nil
		}
	}
	return defValue, nil
}

func (params Params) Set(ctx context.Context, key string, value interface{}) error {
	params[key] = value
	return nil
}

func (params Params) SetInt(ctx context.Context, key string, value int) error {
	params[key] = value
	return nil
}

func (params Params) SetInt8(ctx context.Context, key string, value int8) error {
	params[key] = value
	return nil
}

func (params Params) SetInt16(ctx context.Context, key string, value int16) error {
	params[key] = value
	return nil
}

func (params Params) SetInt32(ctx context.Context, key string, value int32) error {
	params[key] = value
	return nil
}

func (params Params) SetInt64(ctx context.Context, key string, value int64) error {
	params[key] = value
	return nil
}

func (params Params) SetUint(ctx context.Context, key string, value uint) error {
	params[key] = value
	return nil
}

func (params Params) SetUint8(ctx context.Context, key string, value uint8) error {
	params[key] = value
	return nil
}

func (params Params) SetUint16(ctx context.Context, key string, value uint16) error {
	params[key] = value
	return nil
}

func (params Params) SetUint32(ctx context.Context, key string, value uint32) error {
	params[key] = value
	return nil
}

func (params Params) SetUint64(ctx context.Context, key string, value uint64) error {
	params[key] = value
	return nil
}

func (params Params) SetFloat32(ctx context.Context, key string, value float32) error {
	params[key] = value
	return nil
}

func (params Params) SetFloat64(ctx context.Context, key string, value float64) error {
	params[key] = value
	return nil
}

func (params Params) SetString(ctx context.Context, key string, value string) error {
	params[key] = value
	return nil
}

func (params Params) SetBoolean(ctx context.Context, key string, value bool) error {
	params[key] = value
	return nil
}

func (params Params) SetTime(ctx context.Context, key string, value time.Time) error {
	params[key] = value
	return nil
}

func (params Params) JSON() string {
	data, err := json.Marshal(params)
	if err != nil {
		return ""
	}

	return string(data)
}

type reader struct {
	Getter
}

func (reader *reader) Get(ctx context.Context, key string, defValue interface{}) (interface{}, error) {
	value, err := reader.Getter.Get(ctx, key, defValue)
	if err != nil {
		if err == ErrNoMatch {
			return defValue, nil
		}
		return 0, err
	}

	return value, nil
}

func (reader *reader) GetInt(ctx context.Context, key string, defValue int) (int, error) {
	value, err := reader.Getter.Get(ctx, key, defValue)
	if err != nil {
		if err == ErrNoMatch {
			return defValue, nil
		}
		return 0, err
	}

	if val, ok := ConvertToInt(value); ok {
		return val, nil
	}

	return defValue, nil
}

func (reader *reader) GetInt8(ctx context.Context, key string, defValue int8) (int8, error) {
	value, err := reader.Getter.Get(ctx, key, defValue)
	if err != nil {
		if err == ErrNoMatch {
			return defValue, nil
		}
		return 0, err
	}

	if val, ok := ConvertToInt8(value); ok {
		return val, nil
	}

	return defValue, nil
}

func (reader *reader) GetInt16(ctx context.Context, key string, defValue int16) (int16, error) {
	value, err := reader.Getter.Get(ctx, key, defValue)
	if err != nil {
		if err == ErrNoMatch {
			return defValue, nil
		}
		return 0, err
	}

	if val, ok := ConvertToInt16(value); ok {
		return val, nil
	}

	return defValue, nil
}

func (reader *reader) GetInt32(ctx context.Context, key string, defValue int32) (int32, error) {
	value, err := reader.Getter.Get(ctx, key, defValue)
	if err != nil {
		if err == ErrNoMatch {
			return defValue, nil
		}
		return 0, err
	}

	if val, ok := ConvertToInt32(value); ok {
		return val, nil
	}

	return defValue, nil
}

func (reader *reader) GetInt64(ctx context.Context, key string, defValue int64) (int64, error) {
	value, err := reader.Getter.Get(ctx, key, defValue)
	if err != nil {
		if err == ErrNoMatch {
			return defValue, nil
		}
		return 0, err
	}

	if val, ok := ConvertToInt64(value); ok {
		return val, nil
	}

	return defValue, nil
}

func (reader *reader) GetUint(ctx context.Context, key string, defValue uint) (uint, error) {
	value, err := reader.Getter.Get(ctx, key, defValue)
	if err != nil {
		if err == ErrNoMatch {
			return defValue, nil
		}
		return 0, err
	}

	if val, ok := ConvertToUint(value); ok {
		return val, nil
	}

	return defValue, nil
}

func (reader *reader) GetUint8(ctx context.Context, key string, defValue uint8) (uint8, error) {
	value, err := reader.Getter.Get(ctx, key, defValue)
	if err != nil {
		if err == ErrNoMatch {
			return defValue, nil
		}
		return 0, err
	}

	if val, ok := ConvertToUint8(value); ok {
		return val, nil
	}

	return defValue, nil
}

func (reader *reader) GetUint16(ctx context.Context, key string, defValue uint16) (uint16, error) {
	value, err := reader.Getter.Get(ctx, key, defValue)
	if err != nil {
		if err == ErrNoMatch {
			return defValue, nil
		}
		return 0, err
	}

	if val, ok := ConvertToUint16(value); ok {
		return val, nil
	}

	return defValue, nil
}

func (reader *reader) GetUint32(ctx context.Context, key string, defValue uint32) (uint32, error) {
	value, err := reader.Getter.Get(ctx, key, defValue)
	if err != nil {
		if err == ErrNoMatch {
			return defValue, nil
		}
		return 0, err
	}

	if val, ok := ConvertToUint32(value); ok {
		return val, nil
	}

	return defValue, nil
}

func (reader *reader) GetUint64(ctx context.Context, key string, defValue uint64) (uint64, error) {
	value, err := reader.Getter.Get(ctx, key, defValue)
	if err != nil {
		if err == ErrNoMatch {
			return defValue, nil
		}
		return 0, err
	}

	if val, ok := ConvertToUint64(value); ok {
		return val, nil
	}

	return defValue, nil
}

func (reader *reader) GetFloat32(ctx context.Context, key string, defValue float32) (float32, error) {
	value, err := reader.Getter.Get(ctx, key, defValue)
	if err != nil {
		if err == ErrNoMatch {
			return defValue, nil
		}
		return 0, err
	}

	if val, ok := ConvertToFloat32(value); ok {
		return val, nil
	}

	return defValue, nil
}

func (reader *reader) GetFloat64(ctx context.Context, key string, defValue float64) (float64, error) {
	value, err := reader.Getter.Get(ctx, key, defValue)
	if err != nil {
		if err == ErrNoMatch {
			return defValue, nil
		}
		return 0, err
	}

	if val, ok := ConvertToFloat64(value); ok {
		return val, nil
	}

	return defValue, nil
}

func (reader *reader) GetString(ctx context.Context, key string, defValue string) (string, error) {
	value, err := reader.Getter.Get(ctx, key, defValue)
	if err != nil {
		if err == ErrNoMatch {
			return defValue, nil
		}
		return "", err
	}

	if val, ok := ConvertToString(value); ok {
		return val, nil
	}

	return defValue, nil
}

func (reader *reader) GetBoolean(ctx context.Context, key string, defValue bool) (bool, error) {
	value, err := reader.Getter.Get(ctx, key, defValue)
	if err != nil {
		if err == ErrNoMatch {
			return defValue, nil
		}
		return false, err
	}

	if val, ok := ConvertToBoolean(value); ok {
		return val, nil
	}

	return defValue, nil
}

func (reader *reader) GetTime(ctx context.Context, key string, defValue time.Time) (val time.Time, err error) {
	value, err := reader.Getter.Get(ctx, key, defValue)
	if err != nil {
		if err == ErrNoMatch {
			return defValue, nil
		}
		return
	}

	if val, ok := ConvertToTime(value); ok {
		return val, nil
	}

	return defValue, nil
}

type writer struct {
	Setter
}

func (writer *writer) SetInt(ctx context.Context, key string, value int) error {
	return writer.Set(ctx, key, value)
}

func (writer *writer) SetInt8(ctx context.Context, key string, value int8) error {
	return writer.Set(ctx, key, value)
}

func (writer *writer) SetInt16(ctx context.Context, key string, value int16) error {
	return writer.Set(ctx, key, value)
}

func (writer *writer) SetInt32(ctx context.Context, key string, value int32) error {
	return writer.Set(ctx, key, value)
}

func (writer *writer) SetInt64(ctx context.Context, key string, value int64) error {
	return writer.Set(ctx, key, value)
}

func (writer *writer) SetUint(ctx context.Context, key string, value uint) error {
	return writer.Set(ctx, key, value)
}

func (writer *writer) SetUint8(ctx context.Context, key string, value uint8) error {
	return writer.Set(ctx, key, value)
}

func (writer *writer) SetUint16(ctx context.Context, key string, value uint16) error {
	return writer.Set(ctx, key, value)
}

func (writer *writer) SetUint32(ctx context.Context, key string, value uint32) error {
	return writer.Set(ctx, key, value)
}

func (writer *writer) SetUint64(ctx context.Context, key string, value uint64) error {
	return writer.Set(ctx, key, value)
}

func (writer *writer) SetFloat32(ctx context.Context, key string, value float32) error {
	return writer.Set(ctx, key, value)
}

func (writer *writer) SetFloat64(ctx context.Context, key string, value float64) error {
	return writer.Set(ctx, key, value)
}

func (writer *writer) SetString(ctx context.Context, key string, value string) error {
	return writer.Set(ctx, key, value)
}

func (writer *writer) SetTime(ctx context.Context, key string, value time.Time) error {
	return writer.Set(ctx, key, value)
}

func (writer *writer) SetBoolean(ctx context.Context, key string, value bool) error {
	return writer.Set(ctx, key, value)
}

type readerWriter struct {
	reader
	writer
}

func NewReader(getter Getter) Reader {
	return &reader{Getter: getter}
}

func NewWriter(setter Setter) Writer {
	return &writer{Setter: setter}
}

func NewReaderWriter(manager GetterSetter) ReaderWriter {
	return &readerWriter{
		reader: reader{Getter: manager},
		writer: writer{Setter: manager},
	}
}
