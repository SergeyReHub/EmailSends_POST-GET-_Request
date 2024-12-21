package models

import (
	_ "github.com/lib/pq"
)

type Smtp_connection struct {
	Smtp_Host string `json:"smtp_Host"`
	Smtp_Port string `json:"smtp_Port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}
