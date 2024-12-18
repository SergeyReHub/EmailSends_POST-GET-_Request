package content

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"module_git/models"
	"module_git/repository"
	"net/http"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

func Post_contents(ctx *gin.Context) {
	var email models.Email //Getting email struct.
	if err := ctx.ShouldBindJSON(&email); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest) // If email stract was not transferred programm will return StatusBadRequest.
		return
	}

	pool, err := repository.Repository_db()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	err = send_Email(email, ctx)
	var status string
	if err != nil {
		status = fmt.Sprintf("failed: %v", err) // Если отправка не удалась, устанавливаем статус с ошибкой
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
	} else {
		status = "sent" // Если отправка успешна, устанавливаем статус "sent"
	}

	result, err := pool.Exec(ctx, "INSERT INTO messages (subject, body, send_to, status) VALUES ($1, $2, $3, $4)", email.Subject, email.Text, email.To, status)
	if err != nil {
		log.Fatalf("Exec failed: %v\n", err)
	}
	defer pool.Close()

	affectedRows := result.RowsAffected()
	fmt.Printf("Inserted %d rows\n", affectedRows)
	pool.Close()
}

func send_Email(email models.Email, ctx *gin.Context) error {
	var smtp_con models.Smtp_connection
	smtp_con.Smtp_Host = os.Getenv("SMTP_HOST")
	smtp_con.Smtp_Port = os.Getenv("SMTP_PORT")
	smtp_con.Username = os.Getenv("USERNAME_GMAIL")
	smtp_con.Password = os.Getenv("PASSWORD")

	conn, err := smtp.Dial(smtp_con.Smtp_Host + ":" + smtp_con.Smtp_Port) // Creating new client
	if err != nil {
		log.Println("Error connecting to SMTP server:", err)
		return err
	}
	defer conn.Close()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtp_con.Smtp_Host,
	}
	if err = conn.StartTLS(tlsConfig); err != nil {
		log.Println("Error starting TLS:", err)
		return err
	}

	auth := smtp.PlainAuth("", smtp_con.Username, smtp_con.Password, smtp_con.Smtp_Host)
	if err = conn.Auth(auth); err != nil {
		log.Println("Error authenticating:", err)
		return err
	}

	// Отправка письма
	if err = conn.Mail(smtp_con.Username); err != nil {
		log.Println("Error setting sender:", err)
		return err
	}

	pool, err := repository.Repository_db() // Open db connection to RECEPIENTS
	if err != nil {
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}

	// Устанавливаем получателей CC
	for _, cc := range email.CC {
		var s string
		if err = conn.Rcpt(cc); err != nil {
			s = fmt.Sprintf("CC %s: ошибка - %v", cc, err)
			log.Println("Error setting CC recipient:", err)
		} else {
			s = fmt.Sprintf("CC %s: успешно отправлено", cc)
		}

		if err != nil {
			log.Println("Error DB table recepient:", err)
			return err
		}

		last_message_id, err := get_last_message_ID(ctx, pool)
		if err != nil {
			log.Println("Error DB get last message ID:", err)
			return err
		}

		_, dbErr := pool.Exec(ctx, "INSERT INTO recepient (email_adress, status, message_id) VALUES ($1, $2, $3)", cc, s, last_message_id)
		if dbErr != nil {
			log.Printf("Ошибка при добавлении CC %s в базу данных: %v", cc, dbErr)
		}
	}
	pool.Close()

	// Получаем поток для отправки данных
	w, err := conn.Data()
	if err != nil {
		log.Println("Error getting data stream:", err)
		return err
	}

	// Форматируем сообщение
	message := format_Email(email)

	// Записываем сообщение в поток
	_, err = w.Write([]byte(message))
	if err != nil {
		log.Println("Error writing message:", err)
		return err
	}

	err = w.Close()
	if err != nil {
		log.Println("Error closing data stream:", err)
		return err
	}

	err = conn.Quit()
	if err != nil {
		log.Println("Error quiting:", err)
		return err
	}

	// Завершаем сессию
	log.Println("Email sent successfully!")
	return nil
}

func get_last_message_ID(ctx context.Context, pool *pgxpool.Pool) (int64, error) {
	var last_message_ID int64

	// Выполняем SQL-запрос для получения максимального id
	err := pool.QueryRow(ctx, "SELECT COALESCE(MAX(message_id), 0) FROM messages").Scan(&last_message_ID)
	if err != nil {
		log.Printf("Failed to get last message ID: %v\n", err)
		return 0, err // Возвращаем 0 и ошибку в случае неудачи
	}
	last_message_ID++
	return last_message_ID, nil // Возвращаем последний идентификатор
}

func format_Email(email models.Email) string {
	// Форматируем сообщение в соответствии с требованиями SMTP
	return fmt.Sprintf("Subject: %s\r\nTo: %s\r\nCc: %s\r\n\r\n%s",
		email.Subject,
		email.To,
		format_CC(email.CC),
		email.Text)
}

func format_CC(cc []string) string {
	if len(cc) == 0 {
		return ""
	}
	return fmt.Sprintf("%s", cc)
}
