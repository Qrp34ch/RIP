package ds

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWTClaims struct {
	jwt.StandardClaims
	UserID   uint      `json:"user_id"`
	UserUUID uuid.UUID `json:"user_uuid"`
	Scopes   []string  `json:"scopes"`
	IsAdmin  bool
}
