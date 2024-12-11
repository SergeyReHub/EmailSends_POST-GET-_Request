package content

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Config struct {
	Server   Server_config `json:"server"`
	Database DB_config     `json:"database"`
}

type Server_config struct {
	Port string `json:"port"`
}

type DB_config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DB_name  string `json:"db_name"`
	Sslmode  string `json:"sslmode"`
}

func Post_contents(ctx *gin.Context) {
	var email Email //Getting email struct.
	if err := ctx.ShouldBindJSON(&email); err != nil {
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest) // If email stract was not transferred programm will return StatusBadRequest.
		return
	}

	conf_all, err := Load_config()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	} // Getting struct from appsettings.json

	ctx.JSON(http.StatusOK, email)
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		conf_all.Database.User,
		conf_all.Database.Password,
		conf_all.Database.Host,
		conf_all.Database.Port,
		conf_all.Database.DB_name,
		conf_all.Database.Sslmode)
	fmt.Println("Connection string:", connStr)
	db, err := sql.Open("postgres", connStr) // Open db connection
	if err != nil {
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}

	last_message_ID, err := get_last_message_ID(db)
	if err != nil {
		log.Printf("Error getting last message ID: %v", err)
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	log.Printf("lastMessageID: %d;", last_message_ID)
	last_message_ID++

	err = send_Email(email, last_message_ID, &conf_all.Database)
	var status string
	if err != nil {
		status = fmt.Sprintf("failed: %v", err) // Если отправка не удалась, устанавливаем статус с ошибкой
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
	} else {
		status = "sent" // Если отправка успешна, устанавливаем статус "sent"
	}

	result, err := db.Exec("INSERT INTO messages (message_id, subject, body, send_to, status) VALUES ($1, $2, $3, $4, $5)", last_message_ID, email.Subject, email.Text, email.To, status)
	if err != nil {
		log.Printf("Error message exec: %v", err)
		ctx.Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	affectedRows, _ := result.RowsAffected()
	fmt.Printf("Updated %d rows\n", affectedRows)
}

func send_Email(email Email, last_message_ID int, db_conf *DB_config) error {
	smtp_Host := "smtp.gmail.com"
	smtp_Port := "587"
	username := "serzh.rybakov.06@gmail.com" // Ваш адрес электронной почты
	password := "dsvi fgqw kyut bmfa"

	conn, err := smtp.Dial(smtp_Host + ":" + smtp_Port) // Creating new clirnt
	if err != nil {
		fmt.Println("Error connecting to SMTP server:", err)
		return err
	}
	defer conn.Close()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtp_Host,
	}
	if err = conn.StartTLS(tlsConfig); err != nil {
		fmt.Println("Error starting TLS:", err)
		return err
	}

	auth := smtp.PlainAuth("", username, password, smtp_Host)
	if err = conn.Auth(auth); err != nil {
		fmt.Println("Error authenticating:", err)
		return err
	}

	// Отправка письма
	if err = conn.Mail(username); err != nil {
		fmt.Println("Error setting sender:", err)
		return err
	}
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		db_conf.User, db_conf.Password, db_conf.Host, db_conf.Port, db_conf.DB_name, db_conf.Sslmode)
	db, err := sql.Open("postgres", connStr) // Open db connection to RECEPIENTS
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
		last_resepient_ID, err := get_last_recepients_ID(db)
		if err != nil {
			fmt.Println("Error DB table recepient:", err)
			return err
		}
		last_resepient_ID++

		log.Printf("lastMessageID: %d; lastRecepientID: %d", last_message_ID, last_resepient_ID)
		_, dbErr := db.Exec(`INSERT INTO recepient (recepient_id, email_adress, status, "message_id") VALUES ($1, $2, $3, $4)`, last_resepient_ID, cc, s, last_message_ID)
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
	message := format_Email(email)

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

func format_Email(email Email) string {
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

func Load_config() (*Config, error) {
	file, err := os.Open("/home/sergey/my_go_project/config/appsettings.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func get_last_message_ID(db *sql.DB) (int, error) {
	var message_ID sql.NullInt64 // Используем sql.NullInt64 для обработки возможного нулевого значения
	query := "SELECT message_id FROM messages ORDER BY message_id DESC LIMIT 1"

	err := db.QueryRow(query).Scan(&message_ID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если нет строк, возвращаем 0 и nil
			return 0, nil
		}
		return 0, err
	}

	if !message_ID.Valid {
		return 0, nil // Если значение нулевое, возвращаем 1
	}

	return int(message_ID.Int64), nil // Возвращаем полученное значение
}

func get_last_recepients_ID(db *sql.DB) (int, error) {
	var message_ID sql.NullInt64 // Используем sql.NullInt64 для обработки возможного нулевого значения
	query := "SELECT recepient_id FROM recepient ORDER BY recepient_id DESC LIMIT 1"

	err := db.QueryRow(query).Scan(&message_ID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если нет строк, возвращаем 0 и nil
			return 0, nil
		}
		return 0, err
	}

	return int(message_ID.Int64), nil // Возвращаем полученное значение
}

func Get_contents(ctx *gin.Context) {
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
