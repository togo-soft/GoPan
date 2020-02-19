package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

func TestJWT_NewToken(t *testing.T) {
	claims := CustomClaims{
		Username: "test",
		Email:    "test@test.com",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 60,
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "GoPan",
			Subject:   "Login",
		},
	}
	ss, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("R29QYW4="))
	if err != nil {
		panic(err)
	}
	fmt.Println(ss)
}

func TestJWT_ParseToken(t *testing.T) {
	key := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QiLCJlbWFpbCI6InRlc3RAdGVzdC5jb20iLCJleHAiOjE1NzcyODI0MDgsImlhdCI6MTU3NzI4MjM0OCwiaXNzIjoiR29QYW4iLCJuYmYiOjE1NzcyODIzNDgsInN1YiI6IkxvZ2luIn0.cKy0RKnVostTYKjRtNJiuPMynV8opXejuFB8r2SkRS4"
	token, err := jwt.ParseWithClaims(key, &CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("R29QYW4="), nil
		})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		fmt.Printf("%v %v", claims.Username, claims.Email)
	} else {
		fmt.Println("err due to:", err)
	}
}
