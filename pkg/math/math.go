package math

import (
	"crypto/rand"
	"fmt"
	"strconv"
)

type Math interface {
	RandomNumber() *MathResponse
	RandomNumberWithLen(len int) *MathResponse
	AlphaNumericRandom(len int) *MathResponse
}

type MathResponse struct {
	Valid bool
	Value interface{}
}

func (mr *MathResponse) Int() int64 {
	strNumb := fmt.Sprintf("%v", mr.Value)
	numb, err := strconv.Atoi(strNumb)
	if err != nil {
		return 0
	}

	return int64(numb)
}

func (mr *MathResponse) String() string {
	return fmt.Sprintf("%v", mr.Value)
}

type math struct{}

func NewMath() Math {
	return &math{}
}

func (m *math) RandomNumber() *MathResponse {
	// random number with 18 bits
	RandomCrypto, err := rand.Prime(rand.Reader, 18)
	if err != nil {
		return &MathResponse{
			Valid: false,
		}
	}

	return &MathResponse{
		Valid: true,
		Value: fmt.Sprintf("%v", RandomCrypto),
	}
}

func (m *math) RandomNumberWithLen(len int) *MathResponse {
	if len == 0 {
		len = 18
	}

	RandomCrypto, err := rand.Prime(rand.Reader, len)
	if err != nil {
		return &MathResponse{
			Valid: false,
		}
	}

	return &MathResponse{
		Valid: true,
		Value: fmt.Sprintf("%v", RandomCrypto),
	}
}

// len value is max length random alphanumeric
// if len 6 will be provided len * 2 length of random alphanumeric
func (m *math) AlphaNumericRandom(len int) *MathResponse {
	if len == 0 {
		len = 6
	}

	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		return &MathResponse{
			Valid: false,
		}
	}

	return &MathResponse{
		Valid: true,
		Value: fmt.Sprintf("%X", b),
	}
}
