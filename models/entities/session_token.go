package entities

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/twsm000/lenslocked/pkg/crypto/rand"
)

const (
	MinBytesPerToken int = 32
	TokenHashSize    int = 64
)

var (
	ErrTokenSizeBelowMinRequired = errors.New("token size below minimum required")
	ErrFailedToDecodeToken       = errors.New("failed to decode token")
)

type SessionToken struct {
	hash  [TokenHashSize]byte
	value string
}

func (st *SessionToken) Update(size int) error {
	token, err := rand.Bytes(size)
	if err != nil {
		return err
	}

	return st.Set(token)
}

func (st *SessionToken) Set(token []byte) error {
	if len(token) < MinBytesPerToken {
		return ErrTokenSizeBelowMinRequired
	}

	st.value = hex.EncodeToString(token)
	st.hash = sha512.Sum512(token)
	return nil
}

func (st *SessionToken) SetFromHex(hexToken string) error {
	token, err := hex.DecodeString(hexToken)
	if err != nil {
		return errors.Join(ErrFailedToDecodeToken, err)
	}
	return st.Set(token)
}

func (st SessionToken) Value() string {
	return st.value
}

func (st SessionToken) Hash() []byte {
	return st.hash[:]
}

func (st *SessionToken) String() string {
	return hiddenHash
}

func (st *SessionToken) Scan(value any) error {
	st.value = ""
	if value == nil {
		st.hash = [TokenHashSize]byte{}
		return nil
	}

	hash, ok := value.([TokenHashSize]byte)
	if !ok {
		return fmt.Errorf("invalid scan type: %T", value)
	}
	st.hash = hash
	return nil
}
