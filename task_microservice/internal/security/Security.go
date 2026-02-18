package security

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type Security struct{}

func NewSecurity() *Security {
	return &Security{}
}

func (s *Security) GetSubFromToken(tokenString *string) (string, error) {
	block, _ := pem.Decode([]byte("-----BEGIN PUBLIC KEY-----\n" + os.Getenv("keycloak_publickey") + "\n-----END PUBLIC KEY-----"))
	if block == nil || block.Type != "PUBLIC KEY" {
		return "", errors.New("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	rsaPub, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("failed to convert x509.ParsePKIXPublicKey to *rsa.PublicKey")
	}

	token, err := jwt.Parse(*tokenString, func(token *jwt.Token) (any, error) {
		return rsaPub, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))
	if err != nil {
		return "", err
	}

	return token.Claims.(jwt.MapClaims)["sub"].(string), nil
}
