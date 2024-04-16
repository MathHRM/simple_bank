package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"github.com/MathHRM/simple_bank/util"

	_"github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal("Deu bosta man, carregou as config nao: ", err)
	}

	testDB, err = sql.Open(config.DBdriver, config.DBsource)

	if err != nil {
		log.Fatal("Deu erro ai patrion, conectou nao: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}