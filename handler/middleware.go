package handler

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"time"
	"context"
	"fmt"
	"net/http"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
)

type Claims struct {
	Phone string `json:"phone"`
	Id int `json:"id"`
	FullName string `json:"full_name"`
	jwt.StandardClaims
}

var (
	ErrNoAuthHeader      = errors.New("Authorization header is missing")
	ErrInvalidAuthHeader = errors.New("Authorization header is malformed")
	ErrClaimsInvalid     = errors.New("Provided claims do not match expected scopes")
)

// GetJWSFromRequest extracts a JWS string from an Authorization: Bearer <jws> header
func GetJWSFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	// Check for the Authorization header.
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}
	// We expect a header value of the form "Bearer <token>", with 1 space after
	// Bearer, per spec.
	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", ErrInvalidAuthHeader
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}

func loadPrivateKey() *rsa.PrivateKey {
	// get private key file
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

func loadPublicKey() *rsa.PublicKey {
	// get public key file
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

func generateToken(user repository.UserRes) (string, error) {
	// load private key
	privateKey := loadPrivateKey()

	claims := &Claims{
		Phone: user.PhoneNumber.ValueOrZero(),
		Id: user.Id,
		FullName: user.FullName.ValueOrZero(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Hour).Unix(),
		},
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func MiddlewareAuth(c context.Context, input *openapi3filter.AuthenticationInput) error {
	ctx := middleware.GetEchoContext(c)

	// validate security scheme
	if input.SecuritySchemeName != "BearerAuth" {
		ctx.JSON(http.StatusUnauthorized, generated.Message{
			Status:  false,
			Message: fmt.Sprintf(ErrorSecurityScheme, input.SecuritySchemeName),
		})
		return fmt.Errorf(ErrorSecurityScheme, input.SecuritySchemeName)
	}

	// get token from header
	tokenString, err := GetJWSFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, generated.Message{
			Status:  false,
			Message: err.Error(),
		})
		return  err
	}

	// load public key
	publicKey := loadPublicKey()

	// parse claim
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, generated.Message{
			Status:  false,
			Message: "Invalid token",
		})
		return err
	}

	// check if token invalid
	if !token.Valid {
		ctx.JSON(http.StatusUnauthorized, generated.Message{
			Status:  false,
			Message:"Invalid token",
		})
		return err
	}

	// get token claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, generated.Message{
			Status:  false,
			Message: "Invalid token claims",
		})
		return err
	}

	ctx.Set("userId", claims.Id)
	return nil
}
