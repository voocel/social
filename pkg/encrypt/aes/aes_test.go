package aes

import (
	"fmt"
	"testing"
)

func TestAes(t *testing.T) {
	key := "01234567890123456789012345678912"
	str := "hello,world"
	encrypt, err := Encrypt([]byte(key), []byte(str))
	if err != nil {
		fmt.Println(err)
		return
	}
	decrypt, err := Decrypt([]byte(key), encrypt)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("str is %s\n", string(decrypt))
}
