package authentication

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthData struct {
	Token *jwt.Token
	Claim Claims
}

type Claims struct {
	Issuer string `json:"iss"`
	jwt.RegisteredClaims
}

func RequestToToken(req *http.Request) (AuthData, error) {
	t := req.Header.Get("Authorization")

	split := strings.Split(t, " ")
	if len(split) > 1 {
		t = split[1]
	}

	var claims Claims
	token, err := jwt.ParseWithClaims(t, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return AuthData{}, err
	}

	return AuthData{Token: token, Claim: claims}, nil
}

func CreateToken(userID int) (string, string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	sda := jwt.RegisteredClaims{Issuer: "stream-dungeon-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add((time.Hour)).UTC()),
		Subject:   strconv.Itoa(userID)}
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, sda)

	sdr := jwt.RegisteredClaims{Issuer: "stream-dungeon-refresh",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(time.Hour * 24 * 60)).UTC()),
		Subject:   strconv.Itoa(userID)}
	rtk := jwt.NewWithClaims(jwt.SigningMethodHS256, sdr)

	tks, err := tk.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	rtks, err := rtk.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	return tks, rtks, nil
}
