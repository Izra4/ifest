package email

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"text/template"
)

func SendDownloadLink(email, name, link string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "aman-in@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Document Download Link")

	templateFile := "helpers/email/template.html"
	parsedTemplate, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Println(err)
		return
	}

	var tplBuffer bytes.Buffer
	if err := parsedTemplate.Execute(&tplBuffer, map[string]string{
		"Name":         name,
		"DownloadLink": link,
	}); err != nil {
		log.Println("Error execute tempalte")
		return
	}

	m.SetBody("text/html", tplBuffer.String())
	logoPath := "helpers/email/logo.png"
	m.Attach(logoPath, gomail.SetHeader(map[string][]string{
		"Content-ID": {"<logo>"},
	}))
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("EMAIL"), os.Getenv("EMAIL_PASS"))
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
