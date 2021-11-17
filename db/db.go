package db

import (
	"fmt"
	"sync"

	fileutils "github.com/fubss/test-delve-bug-1/fileutils"
	rocksdb "github.com/linxGnu/grocksdb"
)

const (
	closed dbState = iota
	opened
	CurrentFormat = "2.0"
)

type dbState int32

type Conf struct {
	DBPath         string
	ExpectedFormat string
}

// DB - a wrapper on an actual store
type DB struct {
	conf    *Conf
	db      *rocksdb.DB
	dbState dbState
	mutex   sync.RWMutex

	readOpts        *rocksdb.ReadOptions
	writeOptsNoSync *rocksdb.WriteOptions
	writeOptsSync   *rocksdb.WriteOptions
}

// CreateDB constructs a `DB`
func CreateDB(conf *Conf) *DB {
	fmt.Printf("RocksDB constructing...\n")
	readOpts := rocksdb.NewDefaultReadOptions()
	writeOptsNoSync := rocksdb.NewDefaultWriteOptions()
	writeOptsSync := rocksdb.NewDefaultWriteOptions()
	writeOptsSync.SetSync(true)
	fmt.Printf("RocksDB constructing successfully finished\n")
	return &DB{
		conf:            conf,
		dbState:         closed,
		readOpts:        readOpts,
		writeOptsNoSync: writeOptsNoSync,
		writeOptsSync:   writeOptsSync,
	}
}

// Open opens the underlying db
func (dbInst *DB) Open() {
	fmt.Printf("Opening DB in %s...\n", dbInst.conf.DBPath)
	dbInst.mutex.Lock()
	defer dbInst.mutex.Unlock()
	if dbInst.dbState == opened {
		return
	}
	dbOpts := rocksdb.NewDefaultOptions()
	dbPath := dbInst.conf.DBPath
	dbOpts.SetCreateIfMissing(true)
	var err error
	if _, err = fileutils.CreateDirIfMissing(dbPath); err != nil {
		panic(fmt.Sprintf("Error creating dir if missing: %s", err))
	}
	if dbInst.db, err = rocksdb.OpenDb(dbOpts, dbPath); err != nil {
		panic(fmt.Sprintf("Error opening rocksdb: %s", err))
	}
	fmt.Printf("DB was successfully opened\n")
	dbInst.dbState = opened
}
