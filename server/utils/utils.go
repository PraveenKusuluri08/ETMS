package utils

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/smtp"
	"os"
	"text/template"
)

type SendEmailTypes struct {
	To        []string `json:"to"`
	GroupName string   `json:"groupname"`
}

var MAX_NUMBERS = 4

func generateOTPCode() int {
	maxLimit := int64(int(math.Pow10(MAX_NUMBERS)) - 1)
	lowLimit := int(math.Pow10(MAX_NUMBERS - 1))
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(maxLimit))

	if err != nil {
		log.Fatal(err.Error())
	}
	randomIntNumber := int(randomNumber.Int64())

	// for to handle the case of the numbers between from the (0,9999)
	if randomIntNumber < lowLimit {
		randomIntNumber += lowLimit
	}
	// if the random number is greater than the 9999. Never come down just in case
	if randomIntNumber > int(maxLimit) {
		randomIntNumber = int(maxLimit)
	}
	return randomIntNumber
}

func SendEmail(info *SendEmailTypes) string {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")

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
	messageStr := fmt.Sprintf("This OTP use this to accept the invitation $s\n", string(generateOTPCode()))
	err = t.Execute(&body, struct {
		Information string
		Message     string
	}{
		Information: "OTP to accept the invitation",
		Message:     messageStr,
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
