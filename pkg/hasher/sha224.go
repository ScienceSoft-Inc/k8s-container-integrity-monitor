package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

type SHA224 struct {
}

func NewSHA224() *SHA224 {
	return &SHA224{}
}

func (a *SHA224) Hash(file io.Reader) (string, error) {
	hash := sha256.New224()

	_, err := io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
