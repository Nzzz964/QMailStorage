package utils

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path"
)

func MakeChunk(filepath string, chunkSize int64) ([]string, error) {
	info, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	parts := []string{}
	bufffer := make([]byte, chunkSize)
	filesize := info.Size()
	filename := info.Name()
	loop := math.Ceil(float64(filesize) / float64(chunkSize))

	for i := 0; i < int(loop); i++ {
		off := int64(i) * chunkSize

		readSize, err := f.ReadAt(bufffer, off)
		if err != nil && err != io.EOF {
			return parts, err
		}

		tmpName := filename + ".part" + fmt.Sprint(i)
		tmpPath := path.Join(os.TempDir(), tmpName)

		err = ioutil.WriteFile(tmpPath, bufffer[:readSize], 0644)
		if err != nil {
			return parts, err
		}

		parts = append(parts, tmpPath)
	}

	return parts, nil
}

func RemoveChunk(parts []string) {
	for _, v := range parts {
		_ = os.Remove(v)
	}
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
