package passwords

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func New(password string) string {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panicln(err)
	}
	return string(encrypted)
}

func Verify(encrypted, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password)) == nil
}
