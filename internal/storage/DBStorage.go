package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
	"github.com/andrey67895/new_test_go_y_practicum/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStorage struct {
	DB *sql.DB
}

var log = logger.Log()

func (db DBStorage) RetrySaveGauge(id string, delta float64) error {
	err := db.SaveGaugeInDB(id, delta)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if FindErrorInPool(pgErr.Code) {
			for i := 1; i <= 5; i = i + 2 {
				timer := time.NewTimer(time.Duration(i) * time.Second)
				t := <-timer.C
				log.Info(t.Local())
				err := db.SaveGaugeInDB(id, delta)
				if errors.As(err, &pgErr) {
					if !FindErrorInPool(pgErr.Code) {
						break
					}
				}
			}
		}
	}
	return err
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

type Metrics struct {
	ID    string
	MType string
	Delta *float64
	Value *int64
}

func (db DBStorage) GetData() (string, error) {
	data := make([]Metrics, 0)

	rows, err := db.DB.QueryContext(context.Background(), "SELECT * from metrics")
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var v Metrics
		err = rows.Scan(&v.ID, &v.MType, &v.Delta, &v.Value)
		if err != nil {
			return "", err
		}
		data = append(data, v)
	}

	err = rows.Err()
	if err != nil {
		return "", err
	}
	var dataString = ""

	for _, v := range data {
		if strings.ToLower(v.MType) == "gauge" {
			dataString = fmt.Sprintf("%s Name: %s. Delta: %f \n", dataString, v.ID, *v.Delta)
		} else {
			dataString = fmt.Sprintf("%s Name: %s. Value: %d \n", dataString, v.ID, *v.Value)
		}
	}
	return dataString, nil
}

func (db DBStorage) GetCounter(id string) (int64, error) {
	row := db.DB.QueryRowContext(context.Background(), "SELECT m.value as count FROM metrics m WHERE id = $1", id)
	var value int64
	err := row.Scan(&value)
	return value, err
}

func (db DBStorage) GetGauge(id string) (float64, error) {
	row := db.DB.QueryRowContext(context.Background(), "SELECT m.delta as count FROM metrics m WHERE id = $1", id)
	var delta float64
	err := row.Scan(&delta)
	return delta, err
}

func (db DBStorage) SaveGaugeInDB(id string, delta float64) error {
	ctx := context.Background()
	_, err := db.DB.ExecContext(ctx, `INSERT INTO metrics(id, type, delta) values ($1,'GAUGE',$2) on conflict (id) do update set delta = $2`, id, delta)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func InitDB() DBStorage {
	db, err := sql.Open("pgx", config.DatabaseDsn)
	if err != nil {
		log.Error(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		log.Error(err.Error())
	}
	dbStorage := DBStorage{DB: db}
	dbStorage.InitTable()
	return dbStorage
}

func (db DBStorage) InitTable() {
	ctx := context.Background()
	_, err := db.DB.ExecContext(ctx, `DROP TABLE IF EXISTS metrics`)
	if err != nil {
		log.Error(err.Error())
	}
	_, err = db.DB.ExecContext(ctx, `CREATE TABLE metrics (
        "id" text primary key,
        "type" text,
        "delta" double precision,
        "value" bigint
      )`)
	if err != nil {
		log.Error(err.Error())
	}
}

func (db DBStorage) RetrySaveCounter(id string, value int64) error {
	err := db.SaveCounterInDB(id, value)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if FindErrorInPool(pgErr.Code) {
			for i := 1; i <= 5; i = i + 2 {
				timer := time.NewTimer(time.Duration(i) * time.Second)
				t := <-timer.C
				log.Info(t.Local())
				err := db.SaveCounterInDB(id, value)
				if errors.As(err, &pgErr) {
					if !FindErrorInPool(pgErr.Code) {
						break
					}
				}
			}
		}
	}
	return err
}

func (db DBStorage) SaveCounterInDB(id string, value int64) error {
	ctx := context.Background()
	_, err := db.DB.ExecContext(ctx, `INSERT INTO metrics(id, type, value) values ($1,'COUNTER',$2) on conflict (id) do update set value = $2`, id, value)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func (db DBStorage) Ping() error {
	return db.DB.Ping()
}
