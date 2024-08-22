package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
)

func getKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey) { // TODO: Store to disk, this is not usable
	if privateKey != nil && publicKey != nil {
		return privateKey, publicKey
	}

	curve := elliptic.P256()

	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err == nil {
		privateKey = key
		publicKey = &key.PublicKey
	}
	return privateKey, publicKey
}

func generateJWTToken(claims TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	privateKey, _ = getKeys()
	return token.SignedString(privateKey)
}

func VerifyJWTToken(tokenString string) (TokenClaims, error) {
	_, publicKey = getKeys()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return TokenClaims{}, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenClaims := TokenClaims{
			Email: claims["email"].(string),
			ID:    claims["id"].(string),
		}
		return tokenClaims, nil
	}
	return TokenClaims{}, errors.New("invalid token")
}
