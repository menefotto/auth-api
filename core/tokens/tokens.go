package tokens

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"time"

	"github.com/auth-api/core/errors"
	"github.com/auth-api/core/settings"
	"github.com/auth-api/core/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/net/xsrftoken"
)

type customClaims struct {
	Custom string
	jwt.StandardClaims
}

func GenerateJwt(data []byte, delta int) string {
	mapped := string(data)

	claims := customClaims{
		mapped,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(delta)).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "waterandboards",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(utils.GetPrivateKey())
	if err != nil {
		return ""
	}

	return tokenString
}

func ClaimsFromJwt(tok string) (*customClaims, error) {

	token, err := jwt.ParseWithClaims(tok, &customClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.WrongSigningMethod
			}
			// bolocks implementation http://stackoverflow.com/questions/28204385/using-jwt-go-library-key-is-invalid-or-invalid-type
			return utils.GetPubblicKey(), nil
		})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok && !token.Valid {
		return nil, errors.NotValid
	}

	return claims, nil
}

func computeHmac() []byte {
	h := hmac.New(sha512.New, []byte(settings.HMAC_SECRET))

	return h.Sum(nil)
}

type CrsfToken struct {
	Token string `json:"crsf"`
}

func GenerateCrsf(email string) ([]byte, error) {
	payload := xsrftoken.Generate(
		settings.CRYPTO_SECRET,
		email,
		settings.CRSF_ACTION_ID,
	)
	csrf, err := json.Marshal(&CrsfToken{payload})
	if err != nil {
		return nil, err
	}

	return csrf, nil
}
