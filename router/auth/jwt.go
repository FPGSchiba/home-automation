package auth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

var (
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
)

const (
	privateKeyPath = "private_key.pem"
	publicKeyPath  = "public_key.pem"
)

func getKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	if privateKey != nil && publicKey != nil {
		return privateKey, publicKey
	}

	privateKey, publicKey = loadKeysFromDisk()
	if privateKey == nil || publicKey == nil {
		privateKey, publicKey = generateAndSaveKeys()
	}

	return privateKey, publicKey
}

func loadKeysFromDisk() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	privateKeyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, nil
	}

	privateKeyBlock, _ := pem.Decode(privateKeyFile)
	if privateKeyBlock == nil || privateKeyBlock.Type != "EC PRIVATE KEY" {
		return nil, nil
	}

	privateKey, err := x509.ParseECPrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, nil
	}

	publicKeyFile, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, nil
	}

	publicKeyBlock, _ := pem.Decode(publicKeyFile)
	if publicKeyBlock == nil || publicKeyBlock.Type != "PUBLIC KEY" {
		return nil, nil
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, nil
	}

	publicKey, ok := publicKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil
	}

	return privateKey, publicKey
}

func generateAndSaveKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	curve := elliptic.P256()
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil
	}

	privateKey = key
	publicKey = &key.PublicKey

	saveKeysToDisk(privateKey, publicKey)

	return privateKey, publicKey
}

func saveKeysToDisk(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) {
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		fmt.Println("Error marshalling private key:", err)
		return
	}

	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		fmt.Println("Error creating private key file:", err)
		return
	}
	defer privateKeyFile.Close()

	privateKeyBlock := pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	err = pem.Encode(privateKeyFile, &privateKeyBlock)
	if err != nil {
		fmt.Println("Error encoding private key:", err)
		return
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Println("Error marshalling public key:", err)
		return
	}

	publicKeyFile, err := os.Create(publicKeyPath)
	if err != nil {
		fmt.Println("Error creating public key file:", err)
		return
	}
	defer publicKeyFile.Close()

	publicKeyBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	err = pem.Encode(publicKeyFile, &publicKeyBlock)
	if err != nil {
		fmt.Println("Error encoding public key:", err)
		return
	}
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
