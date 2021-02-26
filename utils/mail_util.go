package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	gomail "gopkg.in/gomail.v2"
)

type Info struct {
	Code string
}

func (i Info) SendMailRecovery(umail string) {
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
	m.SetHeader("From", "verificacionmail@rfpjforex.com")
	m.SetHeader("To", umail)
	m.SetHeader("Subject", "Verificación de formulario")
	m.SetBody("text/html", result)

	d := gomail.NewDialer("send.one.com", 465, "verificacionmail@rfpjforex.com", "y%wHW#bCMhdN3dq^")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}
}
