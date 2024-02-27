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

// Update possible errors:
//   - rand.ErrFailedToGenerateSlice
//   - rand.ErrInvalidSizeUnexpected
//   - ErrTokenSizeBelowMinRequired
func (st *SessionToken) Update(size int) Error {
	token, err := rand.Bytes(size)
	if err != nil {
		return NewError(err)
	}

	return st.Set(token)
}

// Set possible errors:
//   - ErrTokenSizeBelowMinRequired
func (st *SessionToken) Set(token []byte) Error {
	if len(token) < MinBytesPerToken {
		return NewError(ErrTokenSizeBelowMinRequired)
	}

	st.value = hex.EncodeToString(token)
	st.hash = sha512.Sum512(token)
	return nil
}

// SetFromHex possible errors:
//   - ErrFailedToDecodeToken
//   - ErrTokenSizeBelowMinRequired
func (st *SessionToken) SetFromHex(hexToken string) Error {
	token, err := hex.DecodeString(hexToken)
	if err != nil {
		return NewError(ErrFailedToDecodeToken, err)
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
		return nil
	}

	hash, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid scan type: %T", value)
	}

	if len(hash) != TokenHashSize {
		return ErrInvalidUser
	}

	copy(st.hash[:], hash)
	return nil
}
