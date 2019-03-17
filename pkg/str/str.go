package str

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

func Rand(n int) string {
	if n <= 0 {
		n = 8
	}
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}
