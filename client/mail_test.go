package client

import (
	"os"
	"qmailstorage/utils"
	"testing"
)

func TestMail(t *testing.T) {
	mail := NewMail("admin@example.com", "admin@example.com")
	mail.Subject("This is Subject")
	mail.Text("This is text/plain")
	mail.Attach("../resources/example.txt", "example.txt")
	f, err := os.Create("mail.example.log")
	if err != nil {
		t.Error(err)
	}
	_, err = mail.WriteTo(f)
	if err != nil {
		t.Error(err)
	}

	f.Close()
	sum, err := utils.Sha1File("mail.example.log")
	if err != nil {
		t.Error(err)
	}
	if sum != "7d8a7b8b95a9dc0ffe7d45f7a6dae7f6565a9665" {
		t.Error("sha1 sum is not correct")
	}

	os.Remove("mail.example.log")
}
