package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/vietthangc1/simple_bank/apis"
	db "github.com/vietthangc1/simple_bank/db/sqlc"
	"github.com/vietthangc1/simple_bank/pkg/envx"
)

var (
	dbDriver    = "postgres"
	postgresUri = envx.String("POSTGRES_URI", "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable")
	address     = envx.String("SERVER_HOST", "localhost:8080")
)

func main() {
	conn, err := sql.Open(dbDriver, postgresUri)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	server := apis.NewServer(store)

	err = server.Start(address)
	if err != nil {
		log.Fatal(err)
	}
}
