package hasher

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
)

type SHA1 struct {
}

func NewSHA1() *SHA1 {
	return &SHA1{}
}

func (a *SHA1) Hash(file io.Reader) (string, error) {
	hash := sha1.New()

	_, err := io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
