package helpers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
	"github.com/andrey67895/new_test_go_y_practicum/internal/logger"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var log = logger.Log()
var DB *sql.DB

func InitDB() *sql.DB {
	db, err := sql.Open("pgx", config.DatabaseDsn)
	if err != nil {
		log.Error(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		log.Error(err.Error())
	}
	return db
}

func InitTable() {
	ctx := context.Background()
	_, err := DB.ExecContext(ctx, `DROP TABLE IF EXISTS metrics`)
	if err != nil {
		log.Error(err.Error())
	}
	_, err = DB.ExecContext(ctx, `CREATE TABLE metrics (
        "id" text primary key,
        "type" text,
        "delta" double precision,
        "value" integer
      )`)
	if err != nil {
		log.Error(err.Error())
	}
}

func SaveGaugeInDB(id string, delta float64) error {
	ctx := context.Background()
	_, err := DB.ExecContext(ctx, `INSERT INTO metrics(id, type, delta) values ($1,'GAUGE',$2) on conflict (id) do update set delta = $2`, id, delta)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func RetrySaveGaugeInDB(id string, delta float64) {
	err := SaveGaugeInDB(id, delta)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if FindErrorInPool(pgErr.Code) {
			for i := 1; i <= 5; i = i + 2 {
				timer := time.NewTimer(time.Duration(i) * time.Second)
				t := <-timer.C
				log.Info(t.Local())
				err := SaveGaugeInDB(id, delta)
				if errors.As(err, &pgErr) {
					if !FindErrorInPool(pgErr.Code) {
						break
					}
				}
			}
		}
	}
}

func RetrySaveCounterInDB(id string, value int64) {
	err := SaveCounterInDB(id, value)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if FindErrorInPool(pgErr.Code) {
			for i := 1; i <= 5; i = i + 2 {
				timer := time.NewTimer(time.Duration(i) * time.Second)
				t := <-timer.C
				log.Info(t.Local())
				err := SaveCounterInDB(id, value)
				if errors.As(err, &pgErr) {
					if !FindErrorInPool(pgErr.Code) {
						break
					}
				}
			}
		}
	}
}

func FindErrorInPool(code string) bool {
	dwarfs := []string{
		pgerrcode.ConnectionException,
		pgerrcode.ConnectionFailure,
		pgerrcode.SQLClientUnableToEstablishSQLConnection,
		pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection,
		pgerrcode.TransactionResolutionUnknown,
		pgerrcode.ProtocolViolation,
		pgerrcode.UniqueViolation}
	return slices.Contains(dwarfs, code)
}

func SaveCounterInDB(id string, value int64) error {
	ctx := context.Background()
	_, err := DB.ExecContext(ctx, `INSERT INTO metrics(id, type, value) values ($1,'COUNTER',$2) on conflict (id) do update set value = $2`, id, value)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func GetPing(w http.ResponseWriter, _ *http.Request) {
	if err := DB.Ping(); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка ping DB: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
