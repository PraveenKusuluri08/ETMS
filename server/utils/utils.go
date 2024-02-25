package utils

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
)

type SendEmailTypes struct {
	To        []string `json:"to"`
	GroupName string   `json:"groupname"`
}

func SendEmail(info *SendEmailTypes) string {
	from := "no.reply.etms@gmail.com"
	password := "asom iatq ufbl abxv"

	var err error
	smtpAuth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	// Initialize a new template
	t := template.New("email.html")

	// Parse the HTML template
	t, err = t.ParseFiles("templates/email.html")
	if err != nil {
		fmt.Println("Error parsing template:", err.Error())
		return "template_error"
	}

	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))
	// Execute the template and write its output to the buffer
	err = t.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    "Praveen",
		Message: "This is a test message in a HTML template",
	})
	if err != nil {
		fmt.Println("Error executing template:", err)
		return "template_execution_error"
	}

	// Send the email using the content from the buffer
	err = smtp.SendMail("smtp.gmail.com:587", smtpAuth, from, info.To, body.Bytes())
	if err != nil {
		fmt.Println("Error sending email:", err)
		return "email_send_error"
	}

	fmt.Println("Email sent successfully!")
	return "email_sent"
}
