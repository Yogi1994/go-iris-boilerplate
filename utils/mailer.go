package utils

import (
	"bytes"
	"html/template"
	"log"

	"github.com/kataras/go-mailer"
)

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func (r *Request) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func SendEmail(email_subject string, templateName string, items interface{}, email_to string) bool {
	r := new(Request)
	err := r.parseTemplate(templateName, items)
	if err != nil {
		log.Fatal(err)
	}

	// sender configuration.
	config := mailer.Config{
		Host:     "EMAIL_HOST",
		Username: "EMAIL_USERNAME",
		Password: "EMAIL_PASSWORD",
		FromAddr: "EMAIL_FROM_ADDRESS",
		Port:     00,
		// Enable UseCommand to support sendmail unix command,
		// if this field is true then Host, Username, Password and Port are not required,
		// because these info already exists in your local sendmail configuration.
		//
		// Defaults to false.
		UseCommand: false,
	}

	// initalize a new mail sender service.
	sender := mailer.New(config)

	// the subject/title of the e-mail.
	subject := email_subject

	// the rich message body.
	content := r.body

	// the recipient(s).
	to := []string{email_to}

	// send the e-mail.
	err = sender.Send(subject, content, to...)

	if err != nil {
		println("error while sending the e-mail: " + err.Error())
		return false
	}
	return true
}
