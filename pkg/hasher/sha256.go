package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

type SHA256 struct {
}

func NewSHA256() *SHA256 {
	return &SHA256{}
}

func (a *SHA256) Hash(file io.Reader) (string, error) {
	hash := sha256.New()

	_, err := io.Copy(hash, file)
	if err != nil {
		fmt.Printf("err in hash copy %s", err)
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
