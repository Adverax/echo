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

package sql

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

type mySqlAdater struct{}

func (adapter *mySqlAdater) Driver() string {
	return "mysql"
}

func (adapter *mySqlAdater) IsDeadlock(db DB, err error) bool {
	return mysqlGetErrorCode(err) == "1213"
}

func (adapter *mySqlAdater) MakeConnectionString(dsn *DSN) string {
	host := dsn.Host
	if host == "" {
		host = "127.0.0.1"
	}
	port := dsn.Port
	if dsn.Port == 0 {
		port = 3306
	}

	var params string
	if dsn.Params == nil {
		params = "?multiStatements=1" // todo: extract params
	} else {
		list := make([]string, 0, len(dsn.Params)+1)
		for key, val := range dsn.Params {
			s := key + "=" + val
			list = append(list, s)
		}
		if _, has := dsn.Params["multiStatements"]; !has {
			list = append(list, "multiStatements=1")
		}
		params = "?" + strings.Join(list, "&")
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s%s",
		dsn.Username,
		dsn.Password,
		host,
		port,
		dsn.Database,
		params,
	)
}

func (adapter *mySqlAdater) DatabaseName(
	db DB,
) (name string, err error) {
	const query = "SELECT database()"
	err = db.QueryRow(query).Scan(&name)
	return
}

// Lock database latch with context
func (adapter *mySqlAdater) LockLocal(ctx context.Context, tx Tx, latch string, timeout int) error {
	var res int
	const query = "SELECT GET_LOCK(CONCAT(DATABASE(), '.', ?), ?)"
	err := tx.QueryRowContext(ctx, query, latch, timeout).Scan(&res)
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrCaptureLock
	}
	return nil
}

// Unlock database with context
func (adapter *mySqlAdater) UnlockLocal(ctx context.Context, tx Tx, latch string) error {
	var res NullInt64
	const query = "SELECT RELEASE_LOCK(CONCAT(DATABASE(), '.', ?))"
	err := tx.QueryRowContext(ctx, query, latch).Scan(&res)
	if err != nil {
		return err
	}
	if !res.Valid {
		return ErrReleaseInvalid
	}
	if res.Int64 == 0 {
		return ErrReleaseLock
	}
	return nil
}

// Lock database latch with context
func (adapter *mySqlAdater) LockGlobal(ctx context.Context, tx Tx, latch string, timeout int) error {
	var res int
	const query = "SELECT GET_LOCK(?, ?)"
	err := tx.QueryRowContext(ctx, query, latch, timeout).Scan(&res)
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrCaptureLock
	}
	return nil
}

// Unlock database with context
func (adapter *mySqlAdater) UnlockGlobal(ctx context.Context, tx Tx, latch string) error {
	var res NullInt64
	const query = "SELECT RELEASE_LOCK(?)"
	err := tx.QueryRowContext(ctx, query, latch).Scan(&res)
	if err != nil {
		return err
	}
	if !res.Valid {
		return ErrReleaseInvalid
	}
	if res.Int64 == 0 {
		return ErrReleaseLock
	}
	return nil
}

func init() {
	Register("mysql", &mySqlAdater{})
}

func mysqlGetErrorCode(err error) string {
	if err == nil {
		return ""
	}

	matches := mysqlErrRe.FindStringSubmatch(err.Error())
	if matches == nil {
		return ""
	}

	return matches[1]
}

var mysqlErrRe = regexp.MustCompile(`Error (\d+):`)
