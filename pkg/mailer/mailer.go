package mailer

import (
	"bytes"
	"errors"
	"text/template"

	"gopkg.in/gomail.v2"
)

type Mail interface {
	Send(
		m Mailer,
		recipient string,
		param any,
	) error
}

type mail struct {
	dialer *gomail.Dialer
	from   string
}

func NewMailer(
	port int,
	username,
	from,
	host,
	password string,
) Mail {
	dialer := gomail.NewDialer(
		host,
		port,
		username,
		password,
	)

	return &mail{
		dialer: dialer,
		from:   from,
	}

}

// params must be on format map[string]interface{}
func (m *mail) Send(
	mail Mailer,
	recipient string,
	param any,
) error {
	temp, ok := mapTemplate[mail]
	if !ok {
		return errors.New("missing template")
	}

	subject, ok := mapSubject[mail]
	if !ok {
		return errors.New("missing subject")
	}

	t, err := template.New("deeplink").Parse(temp)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	err = t.Execute(&b, param)
	if err != nil {
		return err
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", m.from)
	mailer.SetHeader("To", recipient)
	mailer.SetHeader("Subject", subject)
	// mailer.SetBody("text/html", finalTemplate)
	mailer.AddAlternative("text/html", b.String())

	return m.dialer.DialAndSend(mailer)
}
