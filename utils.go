package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randGen    = rand.New(randSource)
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateMembershipID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 16

	for {
		b := make([]byte, length)
		for i := range b {
			b[i] = charset[randGen.Intn(len(charset))]
		}
		membershipID := string(b)

		// Check if the generated ID matches the criteria
		if isValidMembershipID(membershipID) {
			return membershipID
		}
	}
}

func isValidMembershipID(id string) bool {
	if len(id) != 16 {
		return false
	}
	for i := 0; i < len(id); i++ {
		if !((id[i] >= 'A' && id[i] <= 'Z') || (id[i] >= '0' && id[i] <= '9')) {
			return false
		}
		if i > 0 && id[i] == id[i-1] {
			return false
		}
	}
	return true
}
