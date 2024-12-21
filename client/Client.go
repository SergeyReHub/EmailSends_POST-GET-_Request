package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"module_git/models"
	"net/http"
)

func send_email(emailRequest models.Email) error {
	// Преобразование структуры в JSON
	jsonData, err := json.Marshal(emailRequest)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	// Отправка POST-запроса
	resp, err := http.Post("http://app:8080/v1/api/emails", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send email, status code: %d", resp.StatusCode)
	}

	return nil
}

func get_email() error {

	// Отправка POST-запроса
	resp, err := http.Get("http://app:8080/v1/api/emails")
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send email, status code: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	emailRequest := models.Email{
		To:      "serzh.rybakov.06@gmail.com",
		Subject: "Test Subject",
		Text:    "This is a test email.",
		CC:      []string{"serzh.rybakov.06@mail.ru", "jopa342@mail.ru"},
	}

	// Отправка email
	if err := send_email(emailRequest); err != nil {
		fmt.Println("Error sending email (POST query):", err)
		return
	}

	//
	if err := get_email(); err != nil {
		fmt.Println("Error getting emails (GET query):", err)
		return
	}
}
