package main

import (
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/odit-bit/indexstore"
	"github.com/odit-bit/indexstore/indexpostgre"
)

func main() {
	//connect to gpostgre
	var dsn = os.Getenv("DSN")
	db, err := connectPG(dsn)
	if err != nil {
		log.Fatal(err)
	}

	indexer, _ := indexpostgre.New(db)
	idxSrv := indexstore.Server{
		Port:    8383,
		Handler: indexer,
	}

	if err := idxSrv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}

func connectPG(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	count := 3
	for {
		err := db.Ping()
		if err == nil {
			log.Println("success connect to postgre")
			break
		}
		if count > 10 {
			return nil, err
		}
		count *= 2
		dur := time.Duration(count) * time.Second
		log.Printf("try reconnect to postgre in %v \n", dur.Abs().Seconds())
		time.Sleep(dur)

	}
	return db, nil
}
