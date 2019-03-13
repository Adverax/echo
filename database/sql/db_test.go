// +build database

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
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
	"testing"
)

func setup() (DB, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	query, err := ioutil.ReadFile(dir + "/../../_fixture/database/database.sql")
	if err != nil {
		return nil, err
	}

	dsn := &DSC{
		Driver: "mysql",
		DSN: []*DSN{
			{
				Host:     "127.0.0.1",
				Username: "root",
				Password: "SqL314LqS",
			},
		},
	}
	db, err := dsn.Open(nil)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(query))
	if err != nil {
		_ = db.Close(context.Background())
		return nil, err
	}

	return db, nil
}

func TestOpen(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(context.Background())

	if err = db.Ping(); err != nil {
		t.Error(err)
	}
	fmt.Println("DATABASE")
}

func TestClose(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Fatal(err)
	}

	if err = db.Close(context.Background()); err != nil {
		t.Fatal(err)
	}

	if err = db.Ping(); err.Error() != "sql: database is closed" {
		t.Errorf("Physical dbs were not closed correctly. Got: %s", err)
	}
}

/*func TestSlave(t *testing.T) {
	db := &database2{}
	last := -1

	err := quick.Check(func(n int) bool {
		index := db.slave(n)
		if n <= 1 {
			return index == 0
		}

		result := index > 0 && index < n && index != last
		last = index

		return result
	}, nil)

	if err != nil {
		t.Error(err)
	}
}*/
