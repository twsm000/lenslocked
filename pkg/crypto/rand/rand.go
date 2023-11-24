package rand

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

var (
	ErrFailedToGenerateSlice = errors.New("failed to generate slice")
	ErrInvalidSizeUnexpected = errors.New("invalid size unexpected")
)

func Bytes(size int) ([]byte, error) {
	b := make([]byte, size)
	n, err := rand.Read(b)
	if err != nil {
		return nil, errors.Join(ErrFailedToGenerateSlice, err)
	}
	if n != len(b) {
		return nil, ErrInvalidSizeUnexpected
	}
	return b, nil
}

func String(byteSize int) (string, error) {
	data, err := Bytes(byteSize)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(data), nil
}
