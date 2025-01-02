package main

import (
	"FileEncryption/filecrypt"
	"bytes"
	"fmt"
	"os"
	"golang.org/x/term"
)

func printHelp() {
	fmt.Println("Encryption:  go run . encrypt /path/to/file ")
	fmt.Println("Decryption:  go run . decrypt /path/to/file ")
}

func getPassword() ([]byte, []byte) {
	fmt.Println("Enter password:")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		panic(fmt.Sprintf("Error in reading password %v", err))
	}
	fmt.Println("Confirm password:")
	password_2, err_2 := term.ReadPassword(int(os.Stdin.Fd()))
	if err_2 != nil {
		panic(fmt.Sprintf("Error in reading password %v", err_2))
	}

	return password, password_2
}

func validatePassword(password []byte, password_2 []byte) bool {
	if bytes.Equal(password, password_2) {
		return true
	} else {
		return false
	}
}

func validateFile(filepath string) bool {
	_, err := os.Stat(filepath)
	if err == nil {
		return true
	} else {
		fmt.Printf("File does not exist: %v", err)
		return false
	}
}

func isEncrypted(filepath string) bool {
	content, err := os.ReadFile(filepath)
	if err != nil{
		panic(err.Error())
	}
	total_char := len(content)
	if total_char == 0{
		fmt.Println("Empty file")
		return false
	}
	encrypted_char := 0
	for _, b := range content{
		if b != '\n' && b != '\r' {
			ascii := int(b)
			if ascii < 32 {
				encrypted_char++
			}
		}
	}

	threshold := float32(0.3)
	ratio := float32(encrypted_char)/float32(total_char)
	if ratio >= threshold{
		return true
	}

	return false

}

func encrypter(filepath string) {
	pass1, pass2 := getPassword()
	if !validatePassword(pass1, pass2) {
		fmt.Println("Passwords don't match")
		return
	}else{
		if !validateFile(filepath){
			return
		}else{
			if isEncrypted(filepath){
				fmt.Println("File already encrypted")
				return
			}else{
				filecrypt.Encrypt(pass1, filepath)
			}
		}
	}
}

func decryptor(filepath string) {
	pass1, pass2 := getPassword()
	if !validatePassword(pass1, pass2) {
		fmt.Println("Passwords don't match")
		return
	}else{
		if !validateFile(filepath){
			return
		}else{
			if !isEncrypted(filepath){
				fmt.Println("File is already decrypted")
				return
			}
			filecrypt.Decrypt(pass1, filepath)
		}
	}

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Insufficient arguments") //go run . encrypt "/path"
	}

	function_call := os.Args[1]
	fmt.Println(function_call)

	switch function_call {
	case "help":
		printHelp()
	case "encrypt":
		encrypter(os.Args[2])
	case "decrypt":
		decryptor(os.Args[2])
	}
}
