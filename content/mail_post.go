package content

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/smtp"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func getLastMessageID(db *sql.DB) (int, error) {
	var messageID sql.NullInt64 // Используем sql.NullInt64 для обработки возможного нулевого значения
	query := "SELECT message_id FROM messages ORDER BY message_id DESC LIMIT 1"

	err := db.QueryRow(query).Scan(&messageID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если нет строк, возвращаем 0 и nil
			return 0, nil
		}
		return 0, err
	}

	if !messageID.Valid {
		return 0, nil // Если значение нулевое, возвращаем 1
	}

	return int(messageID.Int64), nil // Возвращаем полученное значение
}

func getLastRecepientsID(db *sql.DB) (int, error) {
	var messageID sql.NullInt64 // Используем sql.NullInt64 для обработки возможного нулевого значения
	query := "SELECT recepient_id FROM recepient ORDER BY recepient_id DESC LIMIT 1"

	err := db.QueryRow(query).Scan(&messageID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если нет строк, возвращаем 0 и nil
			return 0, nil
		}
		return 0, err
	}

	return int(messageID.Int64), nil // Возвращаем полученное значение
}

func PostContents(ctx *gin.Context) {
	db, err := sql.Open("postgres", "user=sergey password=stalker1234 host=localhost dbname=mydb sslmode=disable") // Open db connection
	if err != nil {
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	lastMessageID, err := getLastMessageID(db)
	if err != nil {
		log.Printf("Error getting last message ID: %v", err)
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.Printf("lastMessageID: %d;", lastMessageID)
	lastMessageID++

	var email Email
	if err := ctx.ShouldBindJSON(&email); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusOK, email)

	err = sendEmail(email, lastMessageID)
	var status string
	if err != nil {
		status = fmt.Sprintf("failed: %v", err) // Если отправка не удалась, устанавливаем статус с ошибкой
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
	} else {
		status = "sent" // Если отправка успешна, устанавливаем статус "sent"
	}

	result, err := db.Exec("INSERT INTO messages (message_id, subject, body, send_to, status) VALUES ($1, $2, $3, $4, $5)", lastMessageID, email.Subject, email.Text, email.To, status)
	if err != nil {
		log.Printf("Error message exec: %v", err)
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	affectedRows, _ := result.RowsAffected()
	fmt.Printf("Updated %d rows\n", affectedRows)
}

func sendEmail(email Email, messageId int) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	username := "serzh.rybakov.06@gmail.com" // Ваш адрес электронной почты
	password := "dsvi fgqw kyut bmfa"

	conn, err := smtp.Dial(smtpHost + ":" + smtpPort) // Creating new clirnt
	if err != nil {
		fmt.Println("Error connecting to SMTP server:", err)
		return err
	}
	defer conn.Close()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}
	if err = conn.StartTLS(tlsConfig); err != nil {
		fmt.Println("Error starting TLS:", err)
		return err
	}

	auth := smtp.PlainAuth("", username, password, smtpHost)
	if err = conn.Auth(auth); err != nil {
		fmt.Println("Error authenticating:", err)
		return err
	}

	// Отправка письма
	if err = conn.Mail(username); err != nil {
		fmt.Println("Error setting sender:", err)
		return err
	}

	db, err := sql.Open("postgres", "user=sergey password=stalker1234 host=localhost dbname=mydb sslmode=disable") // Open db connection to RECEPIENTS
	if err != nil {
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	// Устанавливаем получателей CC
	for _, cc := range email.CC {
		var s string
		if err = conn.Rcpt(cc); err != nil {
			s = fmt.Sprintf("CC %s: ошибка - %v", cc, err)
			fmt.Println("Error setting CC recipient:", err)

		} else {
			s = fmt.Sprintf("CC %s: успешно отправлено", cc)
		}
		resepientId, err := getLastRecepientsID(db)
		if err != nil {
			fmt.Println("Error DB table recepient:", err)
			return err
		}
		resepientId++

		log.Printf("lastMessageID: %d; lastRecepientID: %d", messageId, resepientId)
		_, dbErr := db.Exec(`INSERT INTO recepient (recepient_id, email_adress, status, "message_id") VALUES ($1, $2, $3, $4)`, resepientId, cc, s, messageId)
		if dbErr != nil {
			log.Printf("Ошибка при добавлении CC %s в базу данных: %v", cc, dbErr)
		}
	}
	db.Close()

	// Получаем поток для отправки данных
	w, err := conn.Data()
	if err != nil {
		fmt.Println("Error getting data stream:", err)
		return err
	}

	// Форматируем сообщение
	message := formatEmail(email)

	// Записываем сообщение в поток
	_, err = w.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing message:", err)
		return err
	}

	err = w.Close()
	if err != nil {
		fmt.Println("Error closing data stream:", err)
		return err
	}

	err = conn.Quit()
	if err != nil {
		fmt.Println("Error quiting:", err)
		return err
	}

	// Завершаем сессию
	fmt.Println("Email sent successfully!")
	return nil
}

func formatEmail(email Email) string {
	// Форматируем сообщение в соответствии с требованиями SMTP
	return fmt.Sprintf("Subject: %s\r\nTo: %s\r\nCc: %s\r\n\r\n%s",
		email.Subject,
		email.To,
		formatCC(email.CC),
		email.Text)
}

func formatCC(cc []string) string {
	if len(cc) == 0 {
		return ""
	}
	return fmt.Sprintf("%s", cc)
}

func GetContents(ctx *gin.Context) {
	name, present := ctx.GetQuery("name")
	if !present {
		ctx.Error(fmt.Errorf("name is required"))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	response := make(map[string]string, 0)
	response["name"] = name

	ctx.JSON(http.StatusOK, response)
}
