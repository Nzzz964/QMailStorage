package utils

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math"
	"os"
	"path"
)

func MakeChunk(filepath string, chunkSize int64) ([]string, error) {
	info, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	src, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	parts := []string{}
	filesize := info.Size()
	filename := info.Name()
	loop := math.Ceil(float64(filesize) / float64(chunkSize))

	for i := 0; i < int(loop); i++ {
		off := int64(i) * chunkSize
		src.Seek(off, io.SeekStart)
		tmpName := filename + ".part" + fmt.Sprint(i)
		tmpPath := path.Join(os.TempDir(), tmpName)
		dst, err := os.Create(tmpPath)
		if err != nil {
			return parts, err
		}

		CopyChunk(dst, src, chunkSize, nil)
		parts = append(parts, tmpPath)
		dst.Close()
	}

	return parts, nil
}

func RemoveChunk(parts []string) {
	for _, v := range parts {
		os.Remove(v)
	}
}

func CopyChunk(dst io.Writer, src io.Reader, need int64, buf []byte) (written int64, err error) {
	if buf == nil {
		buf = make([]byte, 32*1024)
	}

	size := len(buf)
	for {
		nr, er := src.Read(buf)

		distance := need - written
		if distance <= 0 {
			break
		}
		if distance < int64(size) {
			nr = int(distance)
		}

		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

func Sha1File(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
