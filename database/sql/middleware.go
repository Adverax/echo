package sql

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const latchTest = "debug.test"

type Informer interface {
	// Profiler secondary information (skip in production)
	Trace(msg interface{})
}

type Activator func(ctx context.Context, dsc DSC) (DB, error)

// Exclusive database used for exclusive capture database.
// It used latch with database name.
// For latch used separated connection? becouse it keep transaction.
type exclusiveDatabase struct {
	DB
	tx     Tx
	dbname string
}

func (wrapper *exclusiveDatabase) Close(
	ctx context.Context,
) error {
	err := wrapper.DB.Close(ctx)
	if err != nil {
		return err
	}

	err = UnlockGlobal(ctx, wrapper.tx, wrapper.dbname)
	if err != nil {
		return err
	}

	return wrapper.tx.Rollback(ctx)
}

// Open exclusive access for required database
// If control is not null, than for latch opens with heartbeard.
func OpenExclusive(
	ctx context.Context,
	timeout int, // seconds
	activator Activator,
) Activator {
	return func(ctx context.Context, dsc DSC) (DB, error) {
		dsn := dsc.Primary()
		dbname := dsn.Database
		dsn.Database = ""

		/*activator = OpenWithHeartbeart(
			ctx,
			15*time.Second,
			nil,
		)*/

		db, err := Open(ctx, dsc, activator)
		if err != nil {
			return nil, err
		}

		defer func() {
			CloseOnError(ctx, db, err)
		}()

		tx, err := db.Begin(ctx, nil)
		if err != nil {
			return nil, err
		}

		err = LockGlobal(ctx, tx, dbname, timeout)
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
	informer Informer
	indent   string
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
	profiler.informer.Trace(msg)
}

// profiler database profilering text of sql queries.
// Usualy it used for debugging and profiling.
type profilerDB struct {
	DB
	profiler
}

func (db *profilerDB) Begin(ctx context.Context, opts *TxOptions) (Tx, error) {
	org := time.Now()
	tx, err := db.DB.Begin(ctx, opts)
	db.finished("START TRANSACTION", nil, org)
	if err != nil {
		return nil, err
	}
	return &profilerTx{Tx: tx, profiler: db}, nil
}

func (db *profilerDB) Prepare(ctx context.Context, query string) (Stmt, error) {
	stmt, err := db.DB.Prepare(ctx, query)
	if err != nil {
		return nil, err
	}
	return &profilerStmt{Stmt: stmt, profiler: db, query: query}, nil
}

func (db *profilerDB) Exec(ctx context.Context, query string, args ...interface{}) (Result, error) {
	org := time.Now()
	res, err := db.DB.Exec(ctx, query, args...)
	db.finished(query, args, org)
	return res, err
}

func (db *profilerDB) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	org := time.Now()
	res, err := db.DB.Query(ctx, query, args...)
	db.finished(query, args, org)
	return res, err
}

func (db *profilerDB) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	org := time.Now()
	res := db.DB.QueryRow(ctx, query, args...)
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

func (t *profilerTx) Begin(ctx context.Context, opts *TxOptions) (Tx, error) {
	org := time.Now()
	res, err := t.Tx.Begin(ctx, opts)
	t.finished(fmt.Sprintf("SAVEPOINT %d", res.Level()), nil, org)
	return &profilerTx{Tx: res, profiler: t.profiler}, err
}

func (t *profilerTx) Commit(ctx context.Context) error {
	org := time.Now()
	err := t.Tx.Commit(ctx)
	t.finished(fmt.Sprintf("COMMIT %d", t.Level()), nil, org)
	return err
}

func (t *profilerTx) Rollback(ctx context.Context) error {
	org := time.Now()
	err := t.Tx.Rollback(ctx)
	t.finished(fmt.Sprintf("ROLLBACK %d", t.Level()), nil, org)
	return err
}

func (t *profilerTx) Exec(ctx context.Context, query string, args ...interface{}) (Result, error) {
	org := time.Now()
	res, err := t.Tx.Exec(ctx, query, args...)
	t.finished(query, args, org)
	return res, err
}

func (t *profilerTx) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	org := time.Now()
	res, err := t.Tx.Query(ctx, query, args...)
	t.finished(query, args, org)
	return res, err
}

func (t *profilerTx) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	org := time.Now()
	res := t.Tx.QueryRow(ctx, query, args...)
	t.finished(query, args, org)
	return res
}

type profilerStmt struct {
	Stmt
	profiler
	query string
}

func (stmt *profilerStmt) Exec(ctx context.Context, args ...interface{}) (Result, error) {
	org := time.Now()
	res, err := stmt.Stmt.Exec(ctx, args...)
	stmt.finished("EXECUTE STATEMENT\n"+stmt.query, args, org)
	return res, err
}

func (stmt *profilerStmt) Query(ctx context.Context, args ...interface{}) (Rows, error) {
	org := time.Now()
	res, err := stmt.Stmt.Query(ctx, args...)
	stmt.finished("QUERY STATEMENT\n"+stmt.query, args, org)
	return res, err
}

func (stmt *profilerStmt) QueryRow(ctx context.Context, args ...interface{}) Row {
	org := time.Now()
	res := stmt.Stmt.QueryRow(ctx, args...)
	stmt.finished("QUERY STATEMENT\n"+stmt.query, args, org)
	return res
}

// Open database with profiler
func OpenWithProfiler(
	informer Informer,
	indent string,
	activator Activator,
) Activator {
	return func(ctx context.Context, dsc DSC) (DB, error) {
		db, err := Open(ctx, dsc, activator)
		if err != nil {
			return nil, err
		}

		return WithProfiler(db, informer, indent), nil
	}
}

// Wrap database profiler
func WithProfiler(
	db DB,
	informer Informer,
	indent string,
) DB {
	return &profilerDB{
		DB:       db,
		profiler: &profilerMngr{informer: informer, indent: indent},
	}
}

// Open complex database with decorators
// Example:
//   db, err := sql.OpenEx(dsc, reports, autostart, OpenWithHeartbeat(10*time.Second, nil))
func Open(
	ctx context.Context,
	dsc DSC,
	activator Activator,
) (db DB, err error) {
	if activator == nil {
		return open(ctx, dsc, dsc.Type)
	}

	db, err = activator(ctx, dsc)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CloseOnError(ctx context.Context, db DB, err error) {
	if err != nil {
		db.Close(ctx)
	} else {
		e := recover()
		if e != nil {
			db.Close(ctx)
			panic(e)
		}
	}
}
