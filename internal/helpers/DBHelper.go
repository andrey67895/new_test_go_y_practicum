package helpers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func InitDB() *sql.DB {
	db, err := sql.Open("pgx", config.DatabaseDsn)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		panic(err)
	}
	return db
}

func InitTable() {
	ctx := context.Background()
	_, err := DB.ExecContext(ctx, `DROP TABLE IF EXISTS metrics`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.ExecContext(ctx, `CREATE TABLE metrics (
        "id" text primary key,
        "type" text,
        "delta" double precision,
        "value" integer
      )`)
	if err != nil {
		log.Fatal(err)
	}
}

func SaveGaugeInDB(id string, delta float64) {
	log.Printf("ID %s\n", id)
	log.Printf("DELTA %f\n", delta)
	ctx := context.Background()
	_, err := DB.ExecContext(ctx, `INSERT INTO metrics(id, type, delta) values ($1,'GAUGE',$2) on conflict (id) do update set delta = $2`, id, delta)
	if err != nil {
		log.Println(err)
	}
}

func SaveCounterInDB(id string, value int64) {
	log.Printf("ID %s\n", id)
	log.Printf("VALUE %d\n", value)
	ctx := context.Background()
	_, err := DB.ExecContext(ctx, `INSERT INTO metrics(id, type, value) values ($1,'COUNTER',$2) on conflict (id) do update set value = $2`, id, value)
	if err != nil {
		log.Println(err)
	}
}

func GetPing(w http.ResponseWriter, _ *http.Request) {
	if err := DB.Ping(); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка ping DB: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
