package hasher

import (
	"io"
)

//go:generate mockgen -source=hasher.go -destination=mocks/mock_hasher.go

type IHasher interface {
	Hash(file io.Reader) (string, error)
}

// NewHashSum takes a hashing algorithm as input and returns a hash sum with other data or an error
func NewHashSum(alg string) (h IHasher, err error) {
	switch alg {
	case "MD5":
		h = NewMD5()
	case "SHA1":
		h = NewSHA1()
	case "SHA224":
		h = NewSHA224()
	case "SHA384":
		h = NewSHA384()
	case "SHA512":
		h = NewSHA512()
	default:
		h = NewSHA256()
	}
	return h, nil
}
