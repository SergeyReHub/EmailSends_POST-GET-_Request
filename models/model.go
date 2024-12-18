package models

import (
	_ "github.com/lib/pq"
)

type Email struct {
	To      string   `json:"to"`                     // Основные получатели
	CC      []string `json:"carbon_copy_recipients"` // Получатели CC
	Subject string   `json:"subject"`                // Тема письма
	Text    string   `json:"body"`                   // Тело письма
}

type DB_config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DB_name  string `json:"db_name"`
}

type Smtp_connection struct {
	Smtp_Host string `json:"smtp_Host"`
	Smtp_Port string `json:"smtp_Port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}
