package utils

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	enums "coolbreez.lk/moderator/internal/constants"
	"github.com/golang-jwt/jwt/v5"
)

type JWTExtractedDetails struct {
	UserID   int64
	UserRole string
}
type JWTUtil struct {
	secret []byte
}

func NewJWTUtil(jwtSecret string) *JWTUtil {
	return &JWTUtil{
		secret: []byte(jwtSecret),
	}
}

func (ju *JWTUtil) GenerateJWTToken(userID int64, role enums.UserRole) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "coolbreez-moderator",
		"exp":  time.Now().Add(time.Hour * 2).Unix(),
		"sub":  fmt.Sprintf("%d", userID),
		"role": role,
	})
	return token.SignedString(ju.secret)
}

func (ju *JWTUtil) tokenValidator(token *jwt.Token) (any, error) {
	_, isOk := token.Method.(*jwt.SigningMethodHMAC)
	if !isOk {
		return nil, errors.New("unexpected signing method")
	}
	return ju.secret, nil
}

func (ju *JWTUtil) VerifyJWTToken(token string) (*JWTExtractedDetails, error) {
	parsedToken, err := jwt.Parse(token, ju.tokenValidator)
	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("claim type mismatch")
	}
	sub := claims["sub"].(string)
	id, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return nil, errors.New("id convertion failed")
	}

	return &JWTExtractedDetails{
		UserID:   id,
		UserRole: claims["role"].(string),
	}, nil
}
