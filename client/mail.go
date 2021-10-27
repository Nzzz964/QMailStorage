package client

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/mail"
	"strings"
)

const boundary string = "x5fvTI9ZR9aZ"

type Mail struct {
	From    mail.Address
	To      mail.Address
	subject string
	body    []string
}

func NewMail(from string, to string) *Mail {
	return &Mail{
		From: mail.Address{Address: from},
		To:   mail.Address{Address: to},
	}
}

func (mail *Mail) Subject(subject string) *Mail {
	mail.subject = subject
	return mail
}

func (mail *Mail) Text(msg string) *Mail {
	message := strings.Join([]string{
		fmt.Sprintf("--%s", boundary),
		"Content-Type: text/plain; charset=\"utf-8\"",
		"Content-Transfer-Encoding: 7bit",
		"",
		msg,
	}, "\r\n")

	mail.body = append(mail.body, message)
	return mail
}

func (mail *Mail) Html(msg string) *Mail {
	message := strings.Join([]string{
		fmt.Sprintf("--%s", boundary),
		"Content-Type: text/html; charset=\"utf-8\"",
		"Content-Transfer-Encoding: 7bit",
		"",
		msg,
	}, "\r\n")
	mail.body = append(mail.body, message)
	return mail
}

func (mail *Mail) Attach(path string, name string) (*Mail, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return mail, err
	}

	message := strings.Join([]string{
		fmt.Sprintf("--%s", boundary),
		"Content-Type: text/plain; charset=\"utf-8\"",
		"Content-Transfer-Encoding: base64",
		"Content-Disposition: attachment; filename=\"" + name + "\"",
		"",
		base64.StdEncoding.EncodeToString(data),
	}, "\r\n")

	mail.body = append(mail.body, message)
	return mail, nil
}

func (mail *Mail) Reset() *Mail {
	mail.body = nil
	return mail
}

func (mail *Mail) Build() string {
	return strings.Join([]string{
		mail.header(),
		strings.Join(mail.body, "\r\n\r\n"),
	}, "\r\n\r\n")
}

func (mail *Mail) header() string {
	return strings.Join([]string{
		fmt.Sprintf("From: %s", mail.From.String()),
		fmt.Sprintf("To: %s", mail.To.String()),
		fmt.Sprintf("Subject: %s", mail.subject),
		"MIME-Version: 1.0",
		fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"", boundary),
	}, "\r\n")
}
