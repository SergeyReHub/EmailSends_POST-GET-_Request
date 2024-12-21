package content

import (
	"fmt"
	"log"
	"module_git/models"
	"module_git/repository"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func Get_contents(ctx *gin.Context) {
	pool, err := repository.Repository_db()
	if err != nil {
		log.Fatalf("DB pool connect error: %s", err)
	}
	rows_messages, err := pool.Query(ctx, "SELECT * FROM messages") // Замените mytable на имя вашей таблицы
	if err != nil {
		log.Fatalf("error query: %s", err)
	}
	defer rows_messages.Close()

	for rows_messages.Next() {
		var email models.Email
		var message_id int
		var status string

		err := rows_messages.Scan(&message_id, &email.Subject, &email.Text, &email.To, &status)
		if err != nil {
			log.Fatalf("error Scan: %s", err)
		}
		fmt.Printf("ID_message: %d, Subject: %s, Body: %s, Email To: %s, Status: %s\n\n", message_id, email.Subject, email.Text, email.To, status)

		fmt.Printf("carbon_copy_recipients:\n")

		rows_recepients, err := pool.Query(ctx, "SELECT * FROM recepient WHERE message_id = $1", message_id)
		if err != nil {
			log.Fatalf("error query recepients: %s", err)
		}
		defer rows_recepients.Close()
		for rows_recepients.Next() {
			var recepient_id int
			var email_adress string
			var status_recepient string
			err := rows_recepients.Scan(&recepient_id, &email_adress, &status_recepient, nil)
			if err != nil {
				log.Fatalf("error Scan: %s", err)
			}
			fmt.Printf("ID_recepient: %d, email_adress: %s, status_recepient: %s\n", recepient_id, email_adress, status_recepient)

		}
		fmt.Printf("----------------------------------------------->\n\n")
	}
}
