package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash := "$2a$10$B1VbwT58iGEG5vddgpP.MOZMZtYOoyTH9JDDXDnT2ILs3CKrFf/AC"
	password := "caijunjie2"
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("密码不匹配：", err)
	} else {
		fmt.Println("密码匹配，hash有效！")
	}
} 