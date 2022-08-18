package utils

import (
	"log"
	"testing"
)

func TestFFFAddressEncode(t *testing.T) {

	log.Println(FFFAddressEncode("0x0000000000000000000000000000000000002000"))

	log.Println(FFFAddressDecode("FFF5tXp4dJJWemxB4RPMmnVAcanbnsz8MH5Vi3Uf1caR8wtiKtW5nUsGmr"))
}
