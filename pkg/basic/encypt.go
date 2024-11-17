package basic

import (
	"crypto/sha256"
	"fmt"
)

func EncryptPassword(data []byte) (res string) {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
