package client

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/smtp"
	"os"
	"path"
	"qmailstorage/utils"

	"github.com/google/uuid"
)

type Sender struct {
	client   *smtp.Client
	target   string
	desc     string
	fileInfo fs.FileInfo
	checksum string
	chunks   []string
}

type mimeText struct {
	Filename      string `json:"filename"`
	Size          int64  `json:"size"`
	Sha1          string `json:"sha1"`
	Chunksha1     string `json:"chunk_sha1"`
	ChunkFilename string `json:"chunk_filename"`
	Curr          int    `json:"sec"`
	Total         int    `json:"total"`
	Desc          string `json:"desc"`
	Guid          string `json:"guid"`
}

func NewQMailSender(client *smtp.Client, target string, chunksize int64) (sender *Sender, err error) {
	info, err := os.Stat(target)
	if err != nil {
		return nil, err
	}

	sha1str, err := utils.Sha1File(target)
	if err != nil {
		return nil, err
	}

	chunks, err := utils.MakeChunk(target, chunksize)
	if err != nil {
		if len(chunks) > 0 {
			utils.RemoveChunk(chunks)
		}
		return nil, err
	}

	return &Sender{
		client:   client,
		target:   target,
		fileInfo: info,
		checksum: sha1str,
		chunks:   chunks,
	}, nil
}

func (s *Sender) Destory() {
	utils.RemoveChunk(s.chunks)
}

func (s *Sender) Description(desc string) *Sender {
	s.desc = desc
	return s
}

func (s *Sender) Send(from string, to string) error {
	size := s.fileInfo.Size()
	filename := s.fileInfo.Name()

	log.Printf("开始发送邮件...\n")
	log.Printf("文件名: %s\n", filename)
	log.Printf("文件大小: %d bytes\n", size)
	log.Printf("文件 sha1: %s\n", s.checksum)
	log.Printf("文件分片数: %d\n", len(s.chunks))

	guid := uuid.NewString()

	text := mimeText{
		Filename: filename,
		Size:     s.fileInfo.Size(),
		Sha1:     s.checksum,
		Total:    len(s.chunks),
		Desc:     s.desc,
		Guid:     guid,
	}

	mime := NewMail(from, to)
	for k, v := range s.chunks {
		sha1str, err := utils.Sha1File(v)
		if err != nil {
			return err
		}
		basename := path.Base(v)

		text.Chunksha1 = sha1str
		text.ChunkFilename = basename
		text.Curr = k

		jsonBytes, err := json.MarshalIndent(&text, "", "    ")
		if err != nil {
			return err
		}
		json := string(jsonBytes)

		subject := fmt.Sprintf("[QMailStorage服务] 文件: %s 分片: %d", filename, k)
		mime.Subject(subject)
		mime.Text(json)
		mime.Attach(v, basename)

		client := s.client
		if err = client.Mail(from); err != nil {
			return err
		}
		if err = client.Rcpt(to); err != nil {
			return err
		}

		err = func() error {
			w, err := client.Data()
			if err != nil {
				return err
			}
			defer w.Close()
			log.Printf("正在发送邮件: %s\n", subject)
			mime.WriteTo(w)
			return nil
		}()
		if err != nil {
			return err
		}
		mime.Reset()
	}
	return nil
}
