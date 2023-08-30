package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/vietthangc1/simple_bank/pkg/envx"
	"github.com/vietthangc1/simple_bank/pkg/randomx"
)

var (
	dbDriver     = "postgres"
	postgresUri  = envx.String("POSTGRES_URI_DEV", "postgresql://root:secret@localhost:5432/simple_bank_dev?sslmode=disable")
	testQueries  *Queries
	testStore    *Store
	randomEntity randomx.Random
)

func TestMain(m *testing.M) {
	testConn, err := sql.Open(dbDriver, postgresUri)
	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(testConn)
	testStore = NewStore(testConn)
	randomEntity = randomx.NewRandom()

	os.Exit(m.Run())
}
