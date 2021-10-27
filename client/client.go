package client

import (
	"crypto/tls"
	"net"
	"net/smtp"
)

func NewClient(server string, username string, password string) (*smtp.Client, error) {
	host, _, _ := net.SplitHostPort(server)
	auth := smtp.PlainAuth(
		"",
		username,
		password,
		host,
	)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	coon, err := tls.Dial("tcp", server, tlsconfig)
	if err != nil {
		return nil, err
	}
	c, err := smtp.NewClient(coon, host)
	if err != nil {
		return nil, err
	}

	if err = c.Auth(auth); err != nil {
		return nil, err
	}

	return c, nil
}
