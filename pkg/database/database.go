package database

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	sqlConn *sql.DB
	stmt    *Stmt
	singleThreadLock *sync.Mutex
}

func NewDatabase(dbName string) (*Database, error) {
	var err error

	// open database
	sqlConn, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	// prepare statement
	stmt, err := NewPreparedStatement(sqlConn)
	if err != nil {
		return nil, err
	}

	// new database
	database := &Database{
		sqlConn: sqlConn,
		stmt:    stmt,
		singleThreadLock: &sync.Mutex{},
	}

	return database, nil
}
