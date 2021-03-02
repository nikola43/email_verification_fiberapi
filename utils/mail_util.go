package utils

import (
	"bytes"
	"fmt"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	gomail "gopkg.in/gomail.v2"
	"os"
)

type Info struct {
	Code string
	Ip   string
}

func (i Info) SendMailRecoveryEs(umail string) {
	t := template.New("verification_email_template.html")

	var err error
	t, err = t.ParseFiles("verification_email_template.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, i); err != nil {
		log.Println(err)
	}

	result := tpl.String()
	m := gomail.NewMessage()
	m.SetHeader("From", GetEnvVariable("EMAIL"))
	m.SetHeader("To", umail)
	m.SetHeader("Subject", "Verificación de formulario")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("send.one.com", 465, GetEnvVariable("EMAIL"), GetEnvVariable("EMAIL_PASSWORD"))

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}
}

func (i Info) SendMailRecoveryEn(umail string) {
	t := template.New("verification_email_template_en.html")

	var err error
	t, err = t.ParseFiles("verification_email_template_en.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, i); err != nil {
		log.Println(err)
	}

	result := tpl.String()
	m := gomail.NewMessage()
	m.SetHeader("From", GetEnvVariable("EMAIL"))
	m.SetHeader("To", umail)
	m.SetHeader("Subject", "Form verification")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("send.one.com", 465, GetEnvVariable("EMAIL"), GetEnvVariable("EMAIL_PASSWORD"))

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}
}

// use godot package to load/read the .env file and
// return the value of the key
func GetEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
