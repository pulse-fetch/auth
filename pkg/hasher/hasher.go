package hash

import (
	"crypto/sha256"
	"fmt"
)

type Sha256 struct {
	salt string
}

func NewSha256(salt string) *Sha256 {
	return &Sha256{salt: salt}
}
func Hash(sha Sha256, str string) string {
	hashString := sha256.New()
	hashString.Write([]byte(sha.salt))
	hashString.Write([]byte(str))
	return fmt.Sprintf("%x", hashString.Sum(nil))
}
