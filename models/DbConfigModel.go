package models

import (
	_ "github.com/lib/pq"
)

type DB_config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DB_name  string `json:"db_name"`
}