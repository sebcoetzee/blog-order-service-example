package application

import (
	"os"

	"github.com/go-pg/pg"
)

var db *pg.DB

// ResolveDB connects to the database if not already connected or returns a
// database connection if already connected.
func ResolveDB() *pg.DB {
	if db != nil {
		return db
	}

	options, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	options.PoolSize = 5
	db = pg.Connect(options)
	return db
}

// CloseDB closes the database connection.
func CloseDB() {
	if db == nil {
		return
	}

	db.Close()
}
