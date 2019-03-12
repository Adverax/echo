package sql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type DS struct {
	Host     string `json:"host"`     // Host address
	Port     uint16 `json:"port"`     // Host port
	Database string `json:"database"` // Database name
	Username string `json:"username"` // User name
	Password string `json:"password"` // User password
}

// Single dataSource node
type DSN struct {
	Host     string            `json:"host"`     // Host address
	Port     uint16            `json:"port"`     // Host port
	Database string            `json:"database"` // Database name
	Username string            `json:"username"` // User name
	Password string            `json:"password"` // User password
	Params   map[string]string `json:"params"`   // Other parameters
}

type Adapter interface {
	// Get driver name
	Driver() string
	//Get database name
	DatabaseName(ctx context.Context, db DB) (name string, err error)
	// Make connection string for open database
	MakeConnectionString(dsn *DSN) string
	// Check error for deadlock criteria
	IsDeadlock(db DB, err error) bool
	// Acquire local lock
	LockLocal(ctx context.Context, tx Tx, latch string, timeout int) error
	// Release local lock
	UnlockLocal(ctx context.Context, tx Tx, latch string) error
}

type adapterRegistry map[string]Adapter

func (registry adapterRegistry) find(driver string) (Adapter, error) {
	if adapter, ok := registry[driver]; ok {
		return adapter, nil
	}

	return nil, ErrUnknownDriver
}

var adapters = make(adapterRegistry, 8)

func Register(driver string, adapter Adapter) {
	adapters[driver] = adapter
}

func (dsn *DSN) AddParam(key string, value string) {
	if dsn.Params == nil {
		dsn.Params = map[string]string{
			key: value,
		}
	} else {
		dsn.Params[key] = value
	}
}

func (dsn *DSN) Open(
	ctx context.Context,
	driver string,
	activator Activator,
) (DB, error) {
	return Open(
		ctx,
		DSC{
			Driver: driver,
			DSN:    []*DSN{dsn},
		},
		activator,
	)
}

func (dsn *DSN) openSQL(driver string) (*sql.DB, Adapter, error) {
	adapter, err := adapters.find(driver)
	if err != nil {
		return nil, nil, err
	}

	db, err := sql.Open(driver, adapter.MakeConnectionString(dsn))
	if err != nil {
		return nil, nil, err
	}

	return db, adapter, nil
}

// DataSource nodes cluster (first node is master)
type DSC struct {
	Driver string      `json:"Driver"`
	Type   ReactorType `json:"-"`
	DSN    []*DSN
}

func (dsc DSC) Primary() DSN {
	if len(dsc.DSN) == 0 {
		return DSN{}
	}

	return *dsc.DSN[0]
}

func (dsc DSC) String() (string, error) {
	adapter, err := adapters.find(dsc.Driver)
	if err != nil {
		return "", err
	}

	var list = make([]string, len(dsc.DSN))
	for i, dsn := range dsc.DSN {
		list[i] = adapter.MakeConnectionString(dsn)
	}
	return strings.Join(list, ";"), nil
}

func (dsc DSC) Open(
	ctx context.Context,
	activator Activator,
) (DB, error) {
	return Open(ctx, dsc, activator)
}

// Exclusive open database for escape any concurrency.
func (dsc DSC) OpenForTest(
	ctx context.Context,
) DB {
	dsn := dsc.Primary()
	dsn.Database += "_test"
	db, err := dsn.Open(
		ctx,
		dsc.Driver,
		OpenExclusive(ctx, 0x7ffffff, nil),
	)
	if err != nil {
		panic(err)
	}

	return db
}

// IsolationLevel is the transaction isolation level used in TxOptions.
type IsolationLevel int

// Various isolation levels that drivers may support in BeginTx.
// If a driver does not support a given isolation level an error may be returned.
//
// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
const (
	LevelDefault IsolationLevel = iota
	LevelReadUncommitted
	LevelReadCommitted
	LevelWriteCommitted
	LevelRepeatableRead
	LevelSnapshot
	LevelSerializable
	LevelLinearizable
)

type HalfMetrics struct {
	Count int32 `json:"count"` // Count of executed queries
	Time  int64 `json:"time"`  // Elapsed time (microseconds)
}

// Metrics of database
type Metrics struct {
	Query    HalfMetrics `json:"query"`
	Exec     HalfMetrics `json:"exec"`
	Transact HalfMetrics `json:"transact"`
}

func (metrics *Metrics) beginQuery() int64 {
	return time.Now().UnixNano()
}

func (metrics *Metrics) endQuery(started int64) {
	atomic.AddInt32(&metrics.Query.Count, 1)
	atomic.AddInt64(&metrics.Query.Time, time.Now().UnixNano()-started)
}

func (metrics *Metrics) beginExec() int64 {
	return time.Now().UnixNano()
}

func (metrics *Metrics) endExec(started int64) {
	atomic.AddInt32(&metrics.Exec.Count, 1)
	atomic.AddInt64(&metrics.Exec.Time, time.Now().UnixNano()-started)
}

func (metrics *Metrics) beginTransact() int64 {
	return time.Now().UnixNano()
}

func (metrics *Metrics) endTransact(started int64) {
	atomic.AddInt32(&metrics.Transact.Count, 1)
	atomic.AddInt64(&metrics.Transact.Time, time.Now().UnixNano()-started)
}

func (metrics *Metrics) Metrics() Metrics {
	return Metrics{
		Query: HalfMetrics{
			Count: atomic.LoadInt32(&metrics.Query.Count),
			Time:  atomic.LoadInt64(&metrics.Query.Time),
		},
		Exec: HalfMetrics{
			Count: atomic.LoadInt32(&metrics.Exec.Count),
			Time:  atomic.LoadInt64(&metrics.Exec.Time),
		},
		Transact: HalfMetrics{
			Count: atomic.LoadInt32(&metrics.Transact.Count),
			Time:  atomic.LoadInt64(&metrics.Transact.Time),
		},
	}
}

func (metrics *Metrics) Audit(
	ctx context.Context,
	auditor interface{},
) error {
	if a, ok := auditor.(Auditor); ok {
		return a.AuditDatabase(metrics.Metrics())
	}
	return nil
}

type Auditor interface {
	AuditDatabase(Metrics Metrics) error
}

// TxOptions holds the transaction options to be used in DB.BeginTx.
type TxOptions = sql.TxOptions

// Reactor is base interface of DB and Tx
type Reactor interface {
	Begin(ctx context.Context, opts *TxOptions) (Tx, error)
	Exec(ctx context.Context, query string, args ...interface{}) (Result, error)
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	Type() ReactorType
	Adapter() Adapter
}

type Composer interface {
	Add(delta int)
	Done()
	Abort() <-chan struct{}
}

type composer struct {
	sync.WaitGroup
	stop chan struct{}
}

// Stop and wait termination all processes
func (composer *composer) Close() {
	close(composer.stop)
	composer.WaitGroup.Wait()
}

// Get chan for stop handling
func (composer *composer) Abort() <-chan struct{} {
	return composer.stop
}

// DB is a logical database with multiple underlying physical databases
// forming a single master multiple slaves topology.
// Reads and writes are automatically directed to the correct physical db.
type DB interface {
	Reactor
	Composer
	Close(ctx context.Context) error
	Driver() driver.Driver
	Ping() error
	Slave() DB
	Master() DB
	Prepare(ctx context.Context, query string) (Stmt, error)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	SetConnMaxLifetime(d time.Duration)
	IsCluster() bool
	ReactorType() ReactorType
	Metrics() Metrics
	Audit(ctx context.Context, auditor interface{}) error
	Interface(detective func(interface{}) interface{}) (interface{}, bool)
	DSC() DSC
	beginQuery() int64
	endQuery(started int64)
	beginExec() int64
	endExec(started int64)
	beginTransact() int64
	endTransact(started int64)
}

// Scanner is an interface used by Scan.
type Scanner interface {
	// Scan assigns a value from a database driver.
	//
	// The src value will be of one of the following types:
	//
	//    int64
	//    float64
	//    bool
	//    []byte
	//    string
	//    time.Time
	//    nil - for NULL values
	//
	// An error should be returned if the value cannot be stored
	// without loss of information.
	Scan(src interface{}) error
}

// Transaction
type Tx interface {
	Reactor
	Level() int16
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Abstract row fetcher
type Fetcher interface {
	Scan(dest ...interface{}) error
}

// Array of abstract data
type Array []interface{}

type Arrays []Array

func (a Array) Scan(dest ...interface{}) error {
	for i, d := range dest {
		err := ConvertAssign(d, a[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// Executor result
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

// Abstract row
type Row interface {
	Scan(dest ...interface{}) error
}

// Record set
type Rows interface {
	Err() error
	Next() bool
	Columns() ([]string, error)
	Close() error
	Scan(dest ...interface{}) error
}

// Named arguments
type Args map[string]interface{}

// Extract args to usage in Query/Exec
func (args *Args) Extract() []sql.NamedArg {
	res := make([]sql.NamedArg, len(*args))
	var i int
	for key, val := range *args {
		res[i] = sql.NamedArg{Name: key, Value: val}
		i++
	}
	return res
}

type row struct {
	db DB
	r  *sql.Row
}

func (row *row) Scan(dest ...interface{}) error {
	err := row.r.Scan(dest...)
	if err != nil {
		return err
	}
	return nil
}

type rows struct {
	db      DB
	rs      *sql.Rows
	started int64
}

func (rows *rows) Scan(dest ...interface{}) error {
	err := rows.rs.Scan(dest...)
	if err != nil {
		return err
	}
	return nil
}

func (rows *rows) Close() error {
	rows.db.endQuery(rows.started)

	err := rows.rs.Close()
	if err != nil {
		return err
	}
	return nil
}

func (rows *rows) Err() error {
	err := rows.rs.Err()
	if err != nil {
		return err
	}
	return nil
}

func (rows *rows) Next() bool {
	return rows.rs.Next()
}

func (rows *rows) Columns() ([]string, error) {
	res, err := rows.rs.Columns()
	if err != nil {
		return nil, err
	}
	return res, err
}

type result struct {
	db  DB
	res sql.Result
}

func (res *result) LastInsertId() (int64, error) {
	r, err := res.res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return r, nil
}

func (res *result) RowsAffected() (int64, error) {
	r, err := res.res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return r, nil
}

// Open concurrently opens each underlying physical db.
// dataSourceNames must be a semi-comma separated list of DSNs with the first
// one being used as the master and the rest as slaves.
func open(
	dsc DSC,
	reactorType ReactorType,
) (DB, error) {
	if reactorType == 0 {
		reactorType = PrimaryReactor
	}

	if len(dsc.DSN) == 1 {
		db, adapter, err := dsc.DSN[0].openSQL(dsc.Driver)
		if err != nil {
			return nil, err
		}
		return &database1{
			db:          db,
			dsc:         dsc,
			adapter:     adapter,
			reactorType: reactorType,
			composer:    composer{stop: make(chan struct{})},
			Metrics:     new(Metrics),
		}, nil
	}

	adapter, err := adapters.find(dsc.Driver)
	if err != nil {
		return nil, err
	}

	db := &database2{
		pdbs:        make([]*sql.DB, len(dsc.DSN)),
		dsc:         dsc,
		reactorType: reactorType,
		adapter:     adapter,
		composer:    composer{stop: make(chan struct{})},
		Metrics:     new(Metrics),
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

// Wrapper for physical database
type database1 struct {
	db          *sql.DB
	dsc         DSC
	reactorType ReactorType
	adapter     Adapter
	composer
	*Metrics
}

func (db *database1) DSC() DSC {
	return db.dsc
}

func (db *database1) Type() ReactorType {
	return db.reactorType
}

func (db *database1) Adapter() Adapter {
	return db.adapter
}

// Close closes all physical databases concurrently, releasing any open resources.
func (db *database1) Close(ctx context.Context) error {
	db.composer.Close()
	return db.db.Close()
}

// Driver returns the physical database's underlying driver.
func (db *database1) Driver() driver.Driver {
	return db.db.Driver()
}

// Begin starts a transaction on the master. The isolation level is dependent on the driver.
func (db *database1) Begin(ctx context.Context, opts *TxOptions) (Tx, error) {
	started := db.beginTransact()

	t, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &tx{db: db, trans: t, started: started}, nil
}

// Exec executes a query without returning any rows.
// The args are for any named parameters in the query.
// Exec uses the master as the underlying physical db.
func (db *database1) Exec(ctx context.Context, query string, args ...interface{}) (Result, error) {
	started := db.beginExec()
	res, err := db.db.ExecContext(ctx, query, args...)
	db.endExec(started)
	if err != nil {
		return nil, err
	}
	return &result{db: db, res: res}, nil
}

// Prepare creates a prepared statement for later queries or executions
// on each physical database, concurrently.
func (db *database1) Prepare(ctx context.Context, query string) (Stmt, error) {
	s, err := db.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return &stmt1{database: db, stmt: s}, nil
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any parameters in the query.
// Query uses a slave as the physical db.
func (db *database1) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	started := db.beginQuery()
	rs, err := db.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &rows{db: db, rs: rs, started: started}, nil
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always return a non-nil value.
// Errors are deferred until Row's Scan method is called.
// QueryRow uses a slave as the physical db.
func (db *database1) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	started := db.beginQuery()
	r := db.db.QueryRowContext(ctx, query, args...)
	db.endQuery(started)
	return &row{db: db, r: r}
}

// Ping verifies if a connection to database is still alive,
// establishing a connection if necessary.
func (db *database1) Ping() error {
	return db.db.Ping()
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool for each underlying physical db.
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns then the
// new MaxIdleConns will be reduced to match the MaxOpenConns limit
// If n <= 0, no idle connections are retained.
func (db *database1) SetMaxIdleConns(n int) {
	db.db.SetMaxIdleConns(n)
}

// SetMaxOpenConns sets the maximum number of open connections
// to each physical database.
// If MaxIdleConns is greater than 0 and the new MaxOpenConns
// is less than MaxIdleConns, then MaxIdleConns will be reduced to match
// the new MaxOpenConns limit. If n <= 0, then there is no limit on the number
// of open connections. The default is 0 (unlimited).
func (db *database1) SetMaxOpenConns(n int) {
	db.db.SetMaxOpenConns(n)
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
// Expired connections may be closed lazily before reuse.
// If d <= 0, connections are reused forever.
func (db *database1) SetConnMaxLifetime(d time.Duration) {
	db.db.SetConnMaxLifetime(d)
}

// Slave returns one of the physical databases which is a slave
func (db *database1) Slave() DB {
	return db
}

// Master returns the master physical database
func (db *database1) Master() DB {
	return db
}

func (db *database1) IsCluster() bool {
	return false
}

func (db *database1) ReactorType() ReactorType {
	return db.reactorType
}

func (db *database1) Interface(
	detective func(interface{}) interface{},
) (interface{}, bool) {
	res := detective(db)
	if res != nil {
		return res, true
	}
	return nil, false
}

// Cluster database
type database2 struct {
	pdbs        []*sql.DB // Physical databases
	dsc         DSC
	count       uint64 // Monotonically incrementing counter on each query
	reactorType ReactorType
	adapter     Adapter
	composer
	*Metrics
}

func (db *database2) DSC() DSC {
	return db.dsc
}

func (db *database2) Type() ReactorType {
	return db.reactorType
}

func (db *database2) Adapter() Adapter {
	return db.adapter
}

// Close closes all physical databases concurrently, releasing any open resources.
func (db *database2) Close(ctx context.Context) error {
	db.composer.Close()
	return scatter(
		len(db.pdbs),
		func(i int) error {
			err := db.pdbs[i].Close()
			return err
		},
	)
}

// Driver returns the physical database's underlying driver.
func (db *database2) Driver() driver.Driver {
	return db.pdbs[0].Driver()
}

// Begin starts a transaction on the master. The isolation level is dependent on the driver.
func (db *database2) Begin(ctx context.Context, opts *TxOptions) (Tx, error) {
	started := db.beginTransact()

	t, err := db.pdbs[0].BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &tx{db: db, trans: t, started: started}, nil
}

// Exec executes a query without returning any rows.
// The args are for any named parameters in the query.
// Exec uses the master as the underlying physical db.
func (db *database2) Exec(ctx context.Context, query string, args ...interface{}) (Result, error) {
	started := db.beginExec()
	res, err := db.pdbs[0].ExecContext(ctx, query, args...)
	db.endExec(started)
	if err != nil {
		return nil, err
	}
	return &result{db: db, res: res}, nil
}

// Ping verifies if a connection to each physical database is still alive,
// establishing a connection if necessary.
func (db *database2) Ping() error {
	return scatter(
		len(db.pdbs),
		func(i int) error {
			err := db.pdbs[i].Ping()
			return err
		},
	)
}

// Prepare creates a prepared statement for later queries or executions
// on each physical database, concurrently.
func (db *database2) Prepare(ctx context.Context, query string) (Stmt, error) {
	stmts := make([]*sql.Stmt, len(db.pdbs))

	err := scatter(
		len(db.pdbs),
		func(i int) (err error) {
			stmts[i], err = db.pdbs[i].PrepareContext(ctx, query)
			return err
		},
	)

	if err != nil {
		return nil, err
	}

	return &stmt2{db: db, stmts: stmts}, nil
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any parameters in the query.
// Query uses a slave as the physical db.
func (db *database2) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	started := db.beginQuery()
	rs, err := db.pdbs[db.slave(len(db.pdbs))].QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &rows{db: db, rs: rs, started: started}, nil
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always return a non-nil value.
// Errors are deferred until Row's Scan method is called.
// QueryRow uses a slave as the physical db.
func (db *database2) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	started := db.beginQuery()
	r := db.pdbs[db.slave(len(db.pdbs))].QueryRowContext(ctx, query, args...)
	db.endQuery(started)
	return &row{db: db, r: r}
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool for each underlying physical db.
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns then the
// new MaxIdleConns will be reduced to match the MaxOpenConns limit
// If n <= 0, no idle connections are retained.
func (db *database2) SetMaxIdleConns(n int) {
	for i := range db.pdbs {
		db.pdbs[i].SetMaxIdleConns(n)
	}
}

// SetMaxOpenConns sets the maximum number of open connections
// to each physical database.
// If MaxIdleConns is greater than 0 and the new MaxOpenConns
// is less than MaxIdleConns, then MaxIdleConns will be reduced to match
// the new MaxOpenConns limit. If n <= 0, then there is no limit on the number
// of open connections. The default is 0 (unlimited).
func (db *database2) SetMaxOpenConns(n int) {
	for i := range db.pdbs {
		db.pdbs[i].SetMaxOpenConns(n)
	}
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
// Expired connections may be closed lazily before reuse.
// If d <= 0, connections are reused forever.
func (db *database2) SetConnMaxLifetime(d time.Duration) {
	for i := range db.pdbs {
		db.pdbs[i].SetConnMaxLifetime(d)
	}
}

// Slave returns one of the physical databases which is a slave
func (db *database2) Slave() DB {
	return &database1{
		db:      db.pdbs[db.slave(len(db.pdbs))],
		Metrics: db.Metrics,
	}
}

// Master returns the master physical database
func (db *database2) Master() DB {
	return &database1{db: db.pdbs[0]}
}

func (db *database2) slave(n int) int {
	if n <= 1 {
		return 0
	}
	return int(1 + (atomic.AddUint64(&db.count, 1) % uint64(n-1)))
}

func (db *database2) IsCluster() bool {
	return true
}

func (db *database2) ReactorType() ReactorType {
	return db.reactorType
}

func (db *database2) Interface(
	detective func(interface{}) interface{},
) (interface{}, bool) {
	res := detective(db)
	if res != nil {
		return res, true
	}
	return nil, false
}

type tx struct {
	db      DB
	trans   *sql.Tx
	level   int16
	started int64
}

func (t *tx) Type() ReactorType {
	return t.db.Type()
}

func (t *tx) Adapter() Adapter {
	return t.db.Adapter()
}

func (t *tx) Level() int16 {
	return t.level
}

func (t *tx) Begin(ctx context.Context, opts *TxOptions) (Tx, error) {
	res := &txx{tx{db: t.db, trans: t.trans, level: t.level + 1}, true}
	query := "SAVEPOINT " + res.getSavePoint()
	_, err := t.Exec(ctx, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *tx) Commit(ctx context.Context) error {
	err := t.trans.Commit()
	t.db.endTransact(t.started)
	if err != nil {
		return err
	}
	return nil
}

func (t *tx) Rollback(ctx context.Context) error {
	err := t.trans.Rollback()
	if err != nil {
		return err
	}
	return nil
}

func (t *tx) Exec(ctx context.Context, query string, args ...interface{}) (Result, error) {
	started := t.db.beginExec()
	res, err := t.trans.ExecContext(ctx, query, args...)
	t.db.endExec(started)
	if err != nil {
		return nil, err
	}
	return &result{db: t.db, res: res}, nil
}

func (t *tx) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	started := t.db.beginQuery()
	rs, err := t.trans.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &rows{db: t.db, rs: rs, started: started}, nil
}

func (t *tx) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	started := t.db.beginQuery()
	r := t.trans.QueryRowContext(ctx, query, args...)
	t.db.endQuery(started)
	return &row{db: t.db, r: r}
}

func (t *tx) getSavePoint() string {
	return fmt.Sprintf("trans%d", t.level)
}

type txx struct {
	tx
	active bool
}

func (t *txx) Type() ReactorType {
	return t.db.Type()
}

func (t *txx) Adapter() Adapter {
	return t.db.Adapter()
}

func (t *txx) Commit(ctx context.Context) error {
	if !t.active {
		return ErrTxDone
	}
	t.active = false
	query := "RELEASE SAVEPOINT " + t.getSavePoint()
	_, err := t.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (t *txx) Rollback(ctx context.Context) error {
	if !t.active {
		return ErrTxDone
	}
	t.active = false
	query := "ROLLBACK TO SAVEPOINT " + t.getSavePoint()
	_, err := t.Exec(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

var (
	ErrNoRows         = sql.ErrNoRows
	ErrTxDone         = sql.ErrTxDone
	ErrUnknownDriver  = errors.New("unknown driver")
	ErrCaptureLock    = errors.New("timeout of latch")
	ErrReleaseLock    = errors.New("can not release lock")
	ErrReleaseInvalid = errors.New("unknown latch or invalid thread")
)

type Repository interface {
	Reactor(ctx context.Context) Reactor
	Database() DB
	Transaction(
		ctx context.Context,
		action func(ctx context.Context, reactor Reactor) error,
	) error
}

type repository struct {
	reactorKey ReactorType
	db         DB
}

func (repository *repository) Database() DB {
	return repository.db
}

func (repository *repository) Reactor(
	ctx context.Context,
) Reactor {
	db := FromContext(ctx, repository.reactorKey)
	if db != nil {
		return db
	}
	return repository.db
}

func (repository *repository) Transaction(
	ctx context.Context,
	action func(ctx context.Context, reactor Reactor) error,
) (err error) {
	for i := 0; i < 100; i++ {
		err = repository.transaction(ctx, action)
		if err == nil {
			return nil
		}

		if !repository.db.Adapter().IsDeadlock(repository.db, err) {
			return err
		}
	}

	return
}

func (repository *repository) transaction(
	ctx context.Context,
	action func(ctx context.Context, reactor Reactor) error,
) error {
	reactor, err := repository.Reactor(ctx).Begin(ctx, nil)
	if err != nil {
		return err
	}
	defer reactor.Rollback(ctx)

	err = action(WithContext(ctx, reactor), reactor)
	if err != nil {
		return err
	}

	return reactor.Commit(ctx)
}

func NewRepository(db DB) Repository {
	return &repository{
		db:         db,
		reactorKey: db.ReactorType(),
	}
}

// todo: Handle failovers on slaves

/*
Позволяет работать с транзакциями (в том числе и вложенными)
Позволяет инкапсулировать контекст исполнения.
Работа с репликацией
Обработка ошибок
Поддержка оберток (профилирование)
Метрики
*/

/*
Реактор заменить на ?
Удалить контекст везде, кроме репозитория для захвата реактора.

  repository.with(ctx)

context

reactor
tests

*/
/*func (dsn *DSN) String(driver string) (string, error) {
	if adapter, ok := adapters[driver]; ok {
		return adapter.MakeConnectionString(dsn), nil
	}
	return "", ErrUnknownDriver
}*/
func NewDSC(
	driver string,
	reactorType ReactorType,
	dbs []DSN,
) DSC {
	var dsns []*DSN
	for _, db := range dbs {
		params := make(map[string]string)
		for key, val := range db.Params {
			params[key] = val
		}
		dsns = append(
			dsns,
			&DSN{
				Host:     db.Host,
				Port:     db.Port,
				Database: db.Database,
				Username: db.Username,
				Password: db.Password,
				Params:   params,
			},
		)
	}

	return DSC{
		Driver: driver,
		Type:   reactorType,
		DSN:    dsns,
	}
}
