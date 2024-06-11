package helpers

import (
	"context"
	"database/sql"
	"fmt"
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

func GetPing(w http.ResponseWriter, _ *http.Request) {
	if err := DB.Ping(); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка ping DB: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
