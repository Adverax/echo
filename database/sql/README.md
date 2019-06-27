## Features
* There is a wrapper over the standard  database/sql library.
* Allows you to work through Scope, while hiding the differences between the database and the transaction.
* Allows you to transparently pass Scope through Context.
* Allows transparent work with replicas.
* Supports metrics.
* Supports nested transactions.
* Automatic processing of deadlocks.
* Works with various databases (just provide the appropriate adapter).

## Usage
```go
package main

import (
	"context"
	"fmt"
    "log"

    "github.com/adverax/echo/database/sql"
    _ "github.com/go-sql-driver/mysql" // Any sql.DB works
)

type MyRepository struct {
	sql.Repository
}

func (repo *MyRepository) Register(
	ctx context.Context,
	a, b int,
) error {
    return repo.Transaction(
    	ctx,
    	func(ctx context.Context)error{
            scope := repo.Scope(ctx)
            // Working with scope
            _, err := scope.Exec("INSERT INTO sometable1 SET somefield = ?", a)
    		if err != nil {
    			return err
    		}
            _, err = scope.Exec("INSERT INTO sometable2 SET somefield = ?", b)
            return err
    	},
    )
}

func main() {
  // The first DSN is assumed to be the master and all
  // other to be slaves
  dsc := &sql.DSC{
  	Driver: "mysql",
  	DSN: []*sql.DSN{
      {
        Host: "127.0.0.1",
        Database: "echo",
        Username: "root",
        Password: "password",
      },
  	},
  }

  // Use real tracer in next sentence
  db, err := dsc.Open(sql.OpenWithProfiler(nil, "", nil))
  if err != nil {
    log.Fatal(err)
  }
  
  if err := db.Ping(); err != nil {
    log.Fatalf("Some physical database is unreachable: %s", err)
  }

  // Read queries are directed to slaves with Query and QueryRow.
  // Always use Query or QueryRow for SELECTS
  // Load distribution is round-robin only for now.
  var count int
  err = db.QueryRow("SELECT COUNT(*) FROM sometable").Scan(&count)
  if err != nil {
    log.Fatal(err)
  }

  // Write queries are directed to the master with Exec.
  // Always use Exec for INSERTS, UPDATES
  _, err = db.Exec("UPDATE sometable SET something = 1")
  if err != nil {
    log.Fatal(err)
  }

  // Prepared statements are aggregates. If any of the underlying
  // physical databases fails to prepare the statement, the call will
  // return an error. On success, if Exec is called, then the
  // master is used, if Query or QueryRow are called, then a slave
  // is used.
  stmt, err := db.Prepare("SELECT * FROM sometable WHERE something = ?")
  if err != nil {
    log.Fatal(err)
  }
  if _, err := stmt.Exec(); err != nil {
  	log.Fatal(err)
  }

  // Transactions always use the master
  tx, err := db.Begin()
  if err != nil {
    log.Fatal(err)
  }
  // Do something transactional ...
  if err = tx.Commit(); err != nil {
    log.Fatal(err)
  }
  
  // Register data in the repository
  r := &MyRepository{Repository: sql.NewRepository(db)}
  err = r.Register(context.Background(), 10, 20)
  if err != nil {
  	log.Fatal(err)
  }

  // If needed, one can access the master or a slave explicitly.
  master, slave := db.Master(), db.Slave()
  fmt.Println(master.Adapter().DatabaseName(master))
  fmt.Println(slave.Adapter().DatabaseName(slave))
}
```

## Todo
* Support other slave load balancing algorithms.
* Support failovers.
