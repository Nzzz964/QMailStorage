package client

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"net/mail"
	"os"
	"strings"
)

const boundary string = "x5fvTI9ZR9aZ"

type Attachment struct {
	Name string
	Path string
}
type Mail struct {
	From    mail.Address
	To      mail.Address
	subject string
	text    string
	attach  []Attachment
}

var MIMEAttachment string

func init() {
	MIMEAttachment = strings.Join([]string{
		fmt.Sprintf("--%s", boundary),
		"Content-Type: text/plain; charset=\"utf-8\"",
		"Content-Transfer-Encoding: base64",
		"Content-Disposition: attachment; filename=\"%s\"",
		"",
	}, "\r\n")
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

	mail.text = message
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
	mail.text = message
	return mail
}

func (mail *Mail) Attach(path string, name string) *Mail {
	mail.attach = append(mail.attach, Attachment{name, path})
	return mail
}

func (mail *Mail) Reset() *Mail {
	mail.attach = nil
	return mail
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

func (mail *Mail) WriteTo(w io.Writer) (n int64, err error) {
	writer := bufio.NewWriter(w)
	nw, err := writer.WriteString(mail.header())
	n += int64(nw)
	if err != nil {
		return n, err
	}
	nw, err = writer.WriteString("\r\n\r\n")
	n += int64(nw)
	if err != nil {
		return n, err
	}
	nw, err = writer.WriteString(mail.text)
	n += int64(nw)
	if err != nil {
		return n, err
	}
	nw, err = writer.WriteString("\r\n\r\n")
	n += int64(nw)
	if err != nil {
		return n, err
	}
	for _, attachment := range mail.attach {
		nw, err = writer.WriteString(fmt.Sprintf(MIMEAttachment, attachment.Name))
		n += int64(nw)
		if err != nil {
			return n, err
		}
		nw, err = writer.WriteString("\r\n")
		n += int64(nw)
		if err != nil {
			return n, err
		}
		f, err := os.Open(attachment.Path)
		if err != nil {
			return n, err
		}
		reader := bufio.NewReader(f)
		enc := base64.StdEncoding
		base64Writer := base64.NewEncoder(enc, writer)
		reader.WriteTo(base64Writer)
		base64Writer.Close()
		f.Close()
		nw, err = writer.WriteString("\r\n\r\n")
		n += int64(nw)
		if err != nil {
			return n, err
		}
	}
	if b := writer.Buffered(); b > 0 {
		writer.Flush()
	}
	return n, nil
}
