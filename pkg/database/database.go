package database

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	sqlConn          *sql.DB
	stmt             *Stmt
	singleThreadLock SingleThreadLock
}

func NewSingleThreadLock(enabled bool) SingleThreadLock {
	return SingleThreadLock{
		lock:    sync.Mutex{},
		enabled: enabled,
	}
}

type SingleThreadLock struct {
	lock    sync.Mutex
	enabled bool
}

func (stl *SingleThreadLock) Lock() {
	if stl.enabled {
		stl.lock.Lock()
	}
}

func (stl *SingleThreadLock) Unlock() {
	if stl.enabled {
		stl.lock.Unlock()
	}
}

func NewDatabase(dbName string, singleThread bool) (*Database, error) {
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
		sqlConn:          sqlConn,
		stmt:             stmt,
		singleThreadLock: NewSingleThreadLock(singleThread),
	}

	return database, nil
}
