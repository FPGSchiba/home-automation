package util

import (
	"github.com/alexedwards/argon2id"
	log "github.com/sirupsen/logrus"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func VerifyPassword(hashedPassword, password string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hashedPassword)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"component": "util",
			"func":      "VerifyPassword",
		}).Error("Failed to verify password")
		return false
	}
	return match
}
