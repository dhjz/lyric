// /qqmusic/decrypter.go
package qqmusic

import (
	"dlrc/utils/goqrcdec"
	"encoding/hex"
)

func DecryptLyrics(encryptedText string) (string, error) {
	return DecryptLyricsByte(HexStringToByteArray(encryptedText))
}

func DecryptLyricsByte(encrypted []byte) (string, error) {
	res, err := goqrcdec.DecodeQRC(encrypted)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func HexStringToByteArray(hexString string) []byte {
	// log.Println("HexStringToByteArray", hexString)
	length := len(hexString)
	byts := make([]byte, length/2)
	for i := 0; i < length; i += 2 {
		b, _ := hex.DecodeString(hexString[i : i+2])
		byts[i/2] = b[0]
	}
	return byts
}
