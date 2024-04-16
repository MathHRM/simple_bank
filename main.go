package main

import (
	"database/sql"
	"log"

	"github.com/MathHRM/simple_bank/api"
	db "github.com/MathHRM/simple_bank/db/sqlc"
	"github.com/MathHRM/simple_bank/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Deu bosta man, carregou as config nao: ", err)
	}

	conn, err := sql.Open(config.DBdriver, config.DBsource)

	if err != nil {
		log.Fatal("Deu erro ai patrion, conectou no sql nao: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Deu merda em rodar a api: ", err)
	}
}