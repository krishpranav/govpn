package cipher

// imports
import (
	"crypto/rc4"
	"log"
)

var _key = []byte("8pUsXuZw4z6B9EhGdKgNjQnjmVsYv2x5")

func GenerateKey(key string) {
	_key = []byte(key)
}

func XOR(src []byte) []byte {
	c, err := rc4.NewCipher(_key)
	if err != nil {
		log.Fatalln(err)
	}
	dst := make([]byte, len(src))
	c.XORKeyStream(dst, src)
	return dst
}