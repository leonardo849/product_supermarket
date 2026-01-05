package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// type IParser interface {
// 	ParseJWT(tokenString string) (*jwt.MapClaims, error)
// }

type Parser struct {
	secretJwt string
}

// type FakeParser struct {

// }

// func (f *FakeParser) ParseJWT(tokenString string) (*jwt.MapClaims, error) {
// 	claims := jwt.MapClaims{
// 		"auth_id": "",
// 	}
// }

func NewParser(secretJwt string) *Parser {
	return  &Parser{
		secretJwt: secretJwt,
	}
}

func (p *Parser)ParseJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return  nil, fmt.Errorf("unexpected signing method")
			} else {
				return []byte(p.secretJwt), nil
			}
		})
		if err != nil || !token.Valid {
			return  nil, err
		}
		claims := token.Claims.(jwt.MapClaims)
		return &claims, nil
}