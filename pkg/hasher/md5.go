package hasher

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

type MD5 struct {
}

func NewMD5() *MD5 {
	return &MD5{}
}

func (a *MD5) Hash(file io.Reader) (string, error) {
	hash := md5.New()

	_, err := io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
