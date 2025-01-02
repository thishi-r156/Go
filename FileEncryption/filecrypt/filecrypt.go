package filecrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"github.com/xdg-go/pbkdf2"
)



func Encrypt(pass []byte, filepath string) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	nonce := make([]byte, 12)
	_, err = rand.Read(nonce)
	if err != nil {
		log.Fatal(err)
	}

	dk := pbkdf2.Key(pass, nonce, 4096, 32, sha1.New)
	block, err := aes.NewCipher(dk)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)

	if err != nil {
		panic(err.Error())
	}

	cipher_text := aesgcm.Seal(nil, nonce, data, nil)
	cipher_text = append(cipher_text, nonce...)

	dstFile, err := os.Create(filepath)
	if err != nil {
		panic(err.Error())
	}

	defer dstFile.Close()

	_, err = dstFile.Write(cipher_text)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("File encryption successfull")

}

func Decrypt(pass []byte, filepath string) {
	data, err := os.ReadFile(filepath)

	if err != nil {
		panic(err.Error())
	}

	salt := data[len(data)-12:]
	str := hex.EncodeToString(salt)
	nonce, err := hex.DecodeString(str)
	
	if err != nil{
		panic(err.Error())
	}

	dk := pbkdf2.Key(pass, nonce, 4096, 32, sha1.New)

	block, err := aes.NewCipher(dk)

	if err != nil{
		panic(err.Error())
	}

	aes_gcm, err := cipher.NewGCM(block)

	if err != nil{
		panic(err.Error())
	}

	decrypted_text, err := aes_gcm.Open(nil, nonce, data[:len(data)-12], nil)

	if err != nil{
		panic(err.Error())
	}

	dstFile, err := os.Create(filepath)
	if err != nil{
		panic(err.Error())
	}

	defer dstFile.Close()

	_, err = dstFile.Write(decrypted_text)

	if err != nil{
		panic(err.Error())
	}

	fmt.Println("File decrypted successfully")

}
