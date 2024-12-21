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