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
	"fmt"
	"strings"
	"time"
)

type Tracer interface {
	// Profiler secondary information (skip in production)
	Trace(msg interface{})
}

// Exclusive database used for exclusive capture database.
// It used latch with database name.
// For latch used separated connection? becouse it keep transaction.
type exclusiveDatabase struct {
	DB
	tx     Tx
	dbname string
}

func (wrapper *exclusiveDatabase) Close(ctx context.Context) error {
	err := wrapper.DB.Close(ctx)
	if err != nil {
		return err
	}

	err = wrapper.DB.Adapter().UnlockGlobal(ctx, wrapper.tx, wrapper.dbname)
	if err != nil {
		return err
	}

	return wrapper.tx.Rollback()
}

// Open exclusive access for required database
// If control is not null, than for latch opens with heartbeard.
func OpenExclusive(
	ctx context.Context,
	timeout int, // seconds
	activator Activator,
) Activator {
	return func(dsc DSC) (DB, error) {
		dsn := dsc.Primary()
		dbname := dsn.Database
		dsn.Database = ""

		/*activator = OpenWithHeartbeart(
			ctx,
			15*time.Second,
			nil,
		)*/

		db, err := Open(dsc, activator)
		if err != nil {
			return nil, err
		}

		defer func() {
			CloseOnError(ctx, db, err)
		}()

		tx, err := db.Begin()
		if err != nil {
			return nil, err
		}

		err = tx.Adapter().LockGlobal(ctx, tx, dbname, timeout)
		if err != nil {
			return nil, err
		}

		return &exclusiveDatabase{
			DB:     db,
			tx:     tx,
			dbname: dbname,
		}, nil
	}
}

type profiler interface {
	finished(query string, args []interface{}, started time.Time)
}

type profilerMngr struct {
	Tracer
	indent string
}

func (profiler *profilerMngr) finished(
	query string,
	args []interface{},
	started time.Time,
) {
	duration := time.Now().Sub(started)
	query = strings.Trim(query, "\n\t\r ")
	var msg string
	if len(args) == 0 {
		msg = fmt.Sprintf("SQL: Elapsed time %s\n%s",
			duration.String(),
			query,
		)
	} else {
		msg = fmt.Sprintf("SQL: Elapsed time %s\n%s\nArgs: %v",
			duration.String(),
			query,
			args,
		)
	}

	msg = strings.Replace(msg, "\n", "\n"+profiler.indent, -1)
	profiler.Trace(msg)
}

// profiler database profilering text of sql queries.
// Usualy it used for debugging and profiling.
type profilerDB struct {
	DB
	profiler
}

func (db *profilerDB) Begin() (Tx, error) {
	org := time.Now()
	tx, err := db.DB.Begin()
	db.finished("START TRANSACTION", nil, org)
	if err != nil {
		return nil, err
	}
	return &profilerTx{Tx: tx, profiler: db}, nil
}

func (db *profilerDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	org := time.Now()
	tx, err := db.DB.BeginTx(ctx, opts)
	db.finished("START TRANSACTION", nil, org)
	if err != nil {
		return nil, err
	}
	return &profilerTx{Tx: tx, profiler: db}, nil
}

func (db *profilerDB) Prepare(query string) (Stmt, error) {
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &profilerStmt{Stmt: stmt, profiler: db, query: query}, nil
}

func (db *profilerDB) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	stmt, err := db.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &profilerStmt{Stmt: stmt, profiler: db, query: query}, nil
}

func (db *profilerDB) Exec(query string, args ...interface{}) (Result, error) {
	org := time.Now()
	res, err := db.DB.Exec(query, args...)
	db.finished(query, args, org)
	return res, err
}

func (db *profilerDB) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	org := time.Now()
	res, err := db.DB.ExecContext(ctx, query, args...)
	db.finished(query, args, org)
	return res, err
}

func (db *profilerDB) Query(query string, args ...interface{}) (Rows, error) {
	org := time.Now()
	res, err := db.DB.Query(query, args...)
	db.finished(query, args, org)
	return res, err
}

func (db *profilerDB) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	org := time.Now()
	res, err := db.DB.QueryContext(ctx, query, args...)
	db.finished(query, args, org)
	return res, err
}

func (db *profilerDB) QueryRow(query string, args ...interface{}) Row {
	org := time.Now()
	res := db.DB.QueryRow(query, args...)
	db.finished(query, args, org)
	return res
}

func (db *profilerDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	org := time.Now()
	res := db.DB.QueryRowContext(ctx, query, args...)
	db.finished(query, args, org)
	return res
}

func (db *profilerDB) Slave() DB {
	return &profilerDB{
		DB:       db.DB.Slave(),
		profiler: db.profiler,
	}
}

func (db *profilerDB) Master() DB {
	return &profilerDB{
		DB:       db.DB.Master(),
		profiler: db.profiler,
	}
}

type profilerTx struct {
	Tx
	profiler
}

func (t *profilerTx) Begin() (Tx, error) {
	org := time.Now()
	res, err := t.Tx.Begin()
	t.finished(fmt.Sprintf("SAVEPOINT %d", res.Level()), nil, org)
	return &profilerTx{Tx: res, profiler: t.profiler}, err
}

func (t *profilerTx) BeginTx(ctx context.Context, opts *TxOptions) (Tx, error) {
	org := time.Now()
	res, err := t.Tx.BeginTx(ctx, opts)
	t.finished(fmt.Sprintf("SAVEPOINT %d", res.Level()), nil, org)
	return &profilerTx{Tx: res, profiler: t.profiler}, err
}

func (t *profilerTx) Commit() error {
	org := time.Now()
	err := t.Tx.Commit()
	t.finished(fmt.Sprintf("COMMIT %d", t.Level()), nil, org)
	return err
}

func (t *profilerTx) Rollback() error {
	org := time.Now()
	err := t.Tx.Rollback()
	t.finished(fmt.Sprintf("ROLLBACK %d", t.Level()), nil, org)
	return err
}

func (t *profilerTx) Exec(query string, args ...interface{}) (Result, error) {
	org := time.Now()
	res, err := t.Tx.Exec(query, args...)
	t.finished(query, args, org)
	return res, err
}

func (t *profilerTx) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	org := time.Now()
	res, err := t.Tx.ExecContext(ctx, query, args...)
	t.finished(query, args, org)
	return res, err
}

func (t *profilerTx) Query(query string, args ...interface{}) (Rows, error) {
	org := time.Now()
	res, err := t.Tx.Query(query, args...)
	t.finished(query, args, org)
	return res, err
}

func (t *profilerTx) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	org := time.Now()
	res, err := t.Tx.QueryContext(ctx, query, args...)
	t.finished(query, args, org)
	return res, err
}

func (t *profilerTx) QueryRow(query string, args ...interface{}) Row {
	org := time.Now()
	res := t.Tx.QueryRow(query, args...)
	t.finished(query, args, org)
	return res
}

func (t *profilerTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	org := time.Now()
	res := t.Tx.QueryRowContext(ctx, query, args...)
	t.finished(query, args, org)
	return res
}

type profilerStmt struct {
	Stmt
	profiler
	query string
}

func (stmt *profilerStmt) Exec(args ...interface{}) (Result, error) {
	org := time.Now()
	res, err := stmt.Stmt.Exec(args...)
	stmt.finished("EXECUTE STATEMENT\n"+stmt.query, args, org)
	return res, err
}

func (stmt *profilerStmt) ExecContext(ctx context.Context, args ...interface{}) (Result, error) {
	org := time.Now()
	res, err := stmt.Stmt.ExecContext(ctx, args...)
	stmt.finished("EXECUTE STATEMENT\n"+stmt.query, args, org)
	return res, err
}

func (stmt *profilerStmt) Query(args ...interface{}) (Rows, error) {
	org := time.Now()
	res, err := stmt.Stmt.Query(args...)
	stmt.finished("QUERY STATEMENT\n"+stmt.query, args, org)
	return res, err
}

func (stmt *profilerStmt) QueryContext(ctx context.Context, args ...interface{}) (Rows, error) {
	org := time.Now()
	res, err := stmt.Stmt.QueryContext(ctx, args...)
	stmt.finished("QUERY STATEMENT\n"+stmt.query, args, org)
	return res, err
}

func (stmt *profilerStmt) QueryRowContext(ctx context.Context, args ...interface{}) Row {
	org := time.Now()
	res := stmt.Stmt.QueryRowContext(ctx, args...)
	stmt.finished("QUERY STATEMENT\n"+stmt.query, args, org)
	return res
}

// Open database with profiler
func OpenWithProfiler(
	tracer Tracer,
	indent string,
	activator Activator,
) Activator {
	return func(dsc DSC) (DB, error) {
		db, err := Open(dsc, activator)
		if err != nil {
			return nil, err
		}

		return WithProfiler(db, tracer, indent), nil
	}
}

// Wrap database profiler
func WithProfiler(
	db DB,
	tracer Tracer,
	indent string,
) DB {
	return &profilerDB{
		DB:       db,
		profiler: &profilerMngr{Tracer: tracer, indent: indent},
	}
}

// Open complex database with decorators
// Example:
//   db, err := sql.OpenEx(dsc, OpenWithHeartbeat(10*time.Second, nil))
func Open(
	dsc DSC,
	activator Activator,
) (db DB, err error) {
	if activator == nil {
		return open(dsc, dsc.DbId)
	}

	db, err = activator(dsc)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Open concurrently opens each underlying physical db.
// dataSourceNames must be a semi-comma separated list of DSNs with the first
// one being used as the master and the rest as slaves.
func open(
	dsc DSC,
	dbId DbId,
) (DB, error) {
	if dbId == 0 {
		dbId = PrimaryDatabase
	}

	if len(dsc.DSN) == 1 {
		db, adapter, err := dsc.DSN[0].openSQL(dsc.Driver)
		if err != nil {
			return nil, err
		}
		return &database1{
			db:       db,
			dsc:      dsc,
			adapter:  adapter,
			dbId:     dbId,
			composer: composer{stop: make(chan struct{})},
			Metrics:  new(Metrics),
		}, nil
	}

	adapter, err := adapters.find(dsc.Driver)
	if err != nil {
		return nil, err
	}

	db := &database2{
		pdbs:     make([]*sql.DB, len(dsc.DSN)),
		dsc:      dsc,
		dbId:     dbId,
		adapter:  adapter,
		composer: composer{stop: make(chan struct{})},
		Metrics:  new(Metrics),
	}

	err = scatter(
		len(db.pdbs),
		func(i int) (err error) {
			db.pdbs[i], _, err = dsc.DSN[0].openSQL(dsc.Driver)
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseOnError(ctx context.Context, db DB, err error) {
	if err != nil {
		_ = db.Close(ctx)
	} else {
		e := recover()
		if e != nil {
			_ = db.Close(ctx)
			panic(e)
		}
	}
}
