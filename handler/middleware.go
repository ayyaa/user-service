package handler

import (
	"crypto/rsa"
	"io/ioutil"
	"log"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Phone string `json:"phone"`
	jwt.StandardClaims
}

func loadPrivateKey() *rsa.PrivateKey {
	privateKeyFile, err := ioutil.ReadFile("private.pem")
	if err != nil {
		log.Fatal("Error reading private key file:", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyFile)
	if err != nil {
		log.Fatal("Error parsing private key:", err)
	}

	return privateKey
}

func LoadPublicKey() *rsa.PublicKey {
	publicKeyFile, err := ioutil.ReadFile("public.pem")
	if err != nil {
		log.Fatal("Error reading public key file:", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyFile)
	if err != nil {
		log.Fatal("Error parsing public key:", err)
	}

	return publicKey
}

func generateToken(phone string) (string, error) {
	privateKey := loadPrivateKey()

	claims := &Claims{
		Phone: phone,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 4500,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
