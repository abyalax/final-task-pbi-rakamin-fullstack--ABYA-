package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("mn1b2vy3u4i5op6a7q8w9kls0dfghertjcxz")

type JWTClaims struct {
	Username string
	jwt.RegisteredClaims
}
