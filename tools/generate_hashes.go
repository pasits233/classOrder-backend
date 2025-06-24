package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func main() {
	passwords := map[string]string{
		"admin": "admin123",
		"coach": "coach123",
	}

	fmt.Println("--- Generated Hashes ---")
	for name, password := range passwords {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to generate hash for %s: %v", name, err)
		}
		fmt.Printf("'%s': '%s',\n", name, string(hash))
	}
	fmt.Println("------------------------")
} 