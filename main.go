package main

import (
	"database/sql"
	"log"

	"github.com/MathHRM/simple_bank/api"
	db "github.com/MathHRM/simple_bank/db/sqlc"

	_ "github.com/lib/pq"
)


const (
	serverAddress = "127.0.0.1:8080"

	dbDriver = "postgres"
	dbSource = "postgresql://root:senha123@localhost:5432/simple_bank?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("Deu erro ai patrion, conectou no sql nao: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Deu merda em rodar a api: ", err)
	}
}