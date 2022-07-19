package hasher

import (
	"crypto/sha512"
	"encoding/hex"
	"io"
)

type SHA384 struct {
}

func NewSHA384() *SHA384 {
	return &SHA384{}
}

func (a *SHA384) Hash(file io.Reader) (string, error) {
	hash := sha512.New384()

	_, err := io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
