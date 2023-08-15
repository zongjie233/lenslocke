package main

import (
	"fmt"
	"github.com/go-mail/mail/v2"
	_ "github.com/go-mail/mail/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	from := "test@test.com"
	to := "hszdq0608@gmail.com"
	subject := "this is test email"
	plaintext := "This is the body of the email"
	html := `<h1> Hello there !</h1><p>this is the email</p>`

	msg := mail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("From", from)
	msg.SetHeader("Subject", subject)

	msg.SetBody("text/plain", plaintext)

	msg.AddAlternative("text/html", html)
	msg.WriteTo(os.Stdout)

	dialer := mail.NewDialer(host, port, username, password)
	err = dialer.DialAndSend(msg)
	if err != nil {
		panic(err)
	}
	fmt.Println("message sent")

}
