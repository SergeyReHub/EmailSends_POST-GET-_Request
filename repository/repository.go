package repository

import (
	"context"
	"fmt"
	"module_git/models"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

func Repository_db() (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Отменяем контекст, когда функция завершится
	var db_conf models.DB_config
	db_conf.Host = os.Getenv("DB_HOST")
	db_conf.Port = os.Getenv("DB_PORT")
	db_conf.User = os.Getenv("DB_USER")
	db_conf.Password = os.Getenv("DB_PASSWORD")
	db_conf.DB_name = os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db_conf.Host,
		db_conf.Port,
		db_conf.User,
		db_conf.Password,
		db_conf.DB_name)
	pool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
