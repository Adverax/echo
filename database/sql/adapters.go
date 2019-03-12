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
	ctx context.Context,
	db DB,
) (name string, err error) {
	const query = "SELECT database()"
	err = db.QueryRow(ctx, query).Scan(&name)
	return
}

// Lock database latch with context
func (adapter *mySqlAdater) LockLocal(ctx context.Context, tx Tx, latch string, timeout int) error {
	var res int
	query := "SELECT GET_LOCK(CONCAT(DATABASE(), '.', ?), ?)"
	err := tx.QueryRow(ctx, query, latch, timeout).Scan(&res)
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
	err := tx.QueryRow(ctx, "SELECT RELEASE_LOCK(CONCAT(DATABASE(), '.', ?))", latch).Scan(&res)
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
