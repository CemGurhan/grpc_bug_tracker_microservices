package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	// "github.com/dgrijalva/jwt-go"
	user "github.com/cemgurhan/auth-microservice/structs"
	"github.com/golang-jwt/jwt/v4"
)

var MySigningKey = []byte(os.Getenv("SECRET_KEY"))

type GoogleUserService struct {
	email_verified bool
	email          string
	audience       string
}

func getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}

// type GoogleClaims struct {
// 	Email         string `json:"email"`
// 	EmailVerified bool   `json:"email_verified"`
// 	jwt.StandardClaims
// }

var GoogleClaims map[string]interface{}

func IsAuthorized(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := authorize(r)
		if err != nil {
			log.Println(r.Context(), "jwt-auth: %v", err)
		}

		if err == nil && u != nil {
			r = r.WithContext(WithUser(r.Context(), u))
			handle.ServeHTTP(w, r)
		}

	})
}

func authorize(r *http.Request) (*user.GoogleUser, error) {

	if r.Header["Token"] != nil {
		return nil, errors.New("No token found in header")
	}
	// function to parse the token string
	token, err := jwt.Parse(
		r.Header["Token"][0],

		func(token *jwt.Token) (interface{}, error) {

			pem, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))

			_, ok := token.Method.(*jwt.SigningMethodRSA) // validate alg claim
			if !ok {
				return nil, fmt.Errorf("unexpected signing method: %q", token.Header["alg"])
			}

			if err != nil {
				return nil, err
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			if err != nil {
				return nil, fmt.Errorf("Key is not a valid ECDSA public key!")
			}
			return key, nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("could not extract claims (%T): %+v", token.Claims, token.Claims)
	}

	audience := "517952092472-duvetghsstc0deut8fvta8b7n2id8dg5.apps.googleusercontent.com" //not the best way

	if claims["aud"].(string) != audience {
		return nil, fmt.Errorf("mismatched audience. aud field %q does not match %q", claims["aud"], audience)

	}

	email := claims["email"].(string)

	return &user.GoogleUser{
		EmailVerified: true,
		Email:         email,
	}, nil

}