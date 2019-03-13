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
	"database/sql"
)

// Stmt is an aggregate prepared statement.
// It holds a prepared statement for each underlying physical db.
type Stmt interface {
	Close() error
	Exec(args ...interface{}) (Result, error)
	ExecContext(ctx context.Context, args ...interface{}) (Result, error)
	Query(args ...interface{}) (Rows, error)
	QueryContext(ctx context.Context, args ...interface{}) (Rows, error)
	QueryRow(args ...interface{}) Row
	QueryRowContext(ctx context.Context, args ...interface{}) Row
}

// Statement for physical database (dbase)
type stmt1 struct {
	database *database1
	stmt     *sql.Stmt
}

// Close closes the statement by concurrently closing all underlying
// statements concurrently, returning the first non nil error.
func (s *stmt1) Close() error {
	err := s.stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

// Exec executes a prepared statement with the given arguments
// and returns a Result summarizing the effect of the statement.
// Exec uses the master as the underlying physical db.
func (s *stmt1) Exec(args ...interface{}) (Result, error) {
	started := s.database.beginExec()
	res, err := s.stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	s.database.endExec(started)
	return &result{db: s.database, res: res}, nil
}

func (s *stmt1) ExecContext(ctx context.Context, args ...interface{}) (Result, error) {
	started := s.database.beginExec()
	res, err := s.stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	s.database.endExec(started)
	return &result{db: s.database, res: res}, nil
}

// Query executes a prepared query statement with the given
// arguments and returns the query results as a *sql.Rows.
// Query uses a slave as the underlying physical db.
func (s *stmt1) Query(args ...interface{}) (Rows, error) {
	started := s.database.beginQuery()
	rs, err := s.stmt.Query(args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRows
		}
		return nil, err
	}
	return &rows{db: s.database, rs: rs, started: started}, nil
}

func (s *stmt1) QueryContext(ctx context.Context, args ...interface{}) (Rows, error) {
	started := s.database.beginQuery()
	rs, err := s.stmt.QueryContext(ctx, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRows
		}
		return nil, err
	}
	return &rows{db: s.database, rs: rs, started: started}, nil
}

// QueryRow executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error
// will be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *sql.Row's Scan scans the first selected row and discards the rest.
// QueryRow uses a slave as the underlying physical db.
func (s *stmt1) QueryRow(args ...interface{}) Row {
	started := s.database.beginQuery()
	r := s.stmt.QueryRow(args...)
	s.database.endQuery(started)
	return &row{db: s.database, r: r}
}

func (s *stmt1) QueryRowContext(ctx context.Context, args ...interface{}) Row {
	started := s.database.beginQuery()
	r := s.stmt.QueryRowContext(ctx, args...)
	s.database.endQuery(started)
	return &row{db: s.database, r: r}
}

type stmt2 struct {
	db    *database2
	stmts []*sql.Stmt
}

// Close closes the statement by concurrently closing all underlying
// statements concurrently, returning the first non nil error.
func (s *stmt2) Close() error {
	err := scatter(len(s.stmts), func(i int) error {
		return s.stmts[i].Close()
	})
	if err != nil {
		return err
	}
	return nil
}

// Exec executes a prepared statement with the given arguments
// and returns a Result summarizing the effect of the statement.
// Exec uses the master as the underlying physical db.
func (s *stmt2) Exec(args ...interface{}) (Result, error) {
	started := s.db.beginExec()
	res, err := s.stmts[0].Exec(args...)
	if err != nil {
		return nil, err
	}
	s.db.endExec(started)
	return &result{db: s.db, res: res}, nil
}

func (s *stmt2) ExecContext(ctx context.Context, args ...interface{}) (Result, error) {
	started := s.db.beginExec()
	res, err := s.stmts[0].ExecContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	s.db.endExec(started)
	return &result{db: s.db, res: res}, nil
}

// Query executes a prepared query statement with the given
// arguments and returns the query results as a *sql.Rows.
// Query uses a slave as the underlying physical db.
func (s *stmt2) Query(args ...interface{}) (Rows, error) {
	started := s.db.beginQuery()
	rs, err := s.stmts[s.db.slave(len(s.db.pdbs))].Query(args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRows
		}
		return nil, err
	}
	return &rows{db: s.db, rs: rs, started: started}, nil
}

func (s *stmt2) QueryContext(ctx context.Context, args ...interface{}) (Rows, error) {
	started := s.db.beginQuery()
	rs, err := s.stmts[s.db.slave(len(s.db.pdbs))].QueryContext(ctx, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoRows
		}
		return nil, err
	}
	return &rows{db: s.db, rs: rs, started: started}, nil
}

// QueryRow executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error
// will be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *sql.Row's Scan scans the first selected row and discards the rest.
// QueryRow uses a slave as the underlying physical db.
func (s *stmt2) QueryRow(args ...interface{}) Row {
	started := s.db.beginQuery()
	r := s.stmts[s.db.slave(len(s.db.pdbs))].QueryRow(args...)
	s.db.endQuery(started)
	return &row{db: s.db, r: r}
}

func (s *stmt2) QueryRowContext(ctx context.Context, args ...interface{}) Row {
	started := s.db.beginQuery()
	r := s.stmts[s.db.slave(len(s.db.pdbs))].QueryRowContext(ctx, args...)
	s.db.endQuery(started)
	return &row{db: s.db, r: r}
}
