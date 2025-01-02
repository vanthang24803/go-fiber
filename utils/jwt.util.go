package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vanthang24803/go-api/infra"
	"github.com/vanthang24803/go-api/internal/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type JwtPayload struct {
	Sub      primitive.ObjectID `json:"sub"`
	Iat      int64              `json:"iat"`
	Exp      int64              `json:"exp"`
	Username string             `json:"username"`
	Roles    []string           `json:"roles"`
}

func GenerateJWT(user *schema.User) (*TokenResponse, error) {
	config := infra.GetConfig()
	now := time.Now()

	// Tạo access token
	rolesStr := strings.Join(user.Roles, ",") // Chuyển đổi mảng Roles thành chuỗi
	accessPayload := JwtPayload{
		Sub:      user.ID,
		Iat:      now.Unix(),
		Exp:      now.Add(time.Hour * 72).Unix(),
		Username: user.Username,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      accessPayload.Sub,
		"iat":      accessPayload.Iat,
		"exp":      accessPayload.Exp,
		"username": accessPayload.Username,
		"roles":    rolesStr, // Nhét roles dưới dạng chuỗi
	})

	accessTokenString, err := accessToken.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return nil, err
	}

	// Tạo refresh token
	refreshPayload := JwtPayload{
		Sub: user.ID,
		Iat: now.Unix(),
		Exp: now.Add(time.Hour * 24 * 7).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": refreshPayload.Sub,
		"iat": refreshPayload.Iat,
		"exp": refreshPayload.Exp,
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(config.JWTRefreshSecret))
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func ValidateJWT(tokenString string) (*JwtPayload, error) {
	config := infra.GetConfig()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		subStr, ok := claims["sub"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid sub claim")
		}

		subStr = strings.Trim(subStr, "ObjectID()")
		subStr = strings.Trim(subStr, "\"")

		subID, err := primitive.ObjectIDFromHex(subStr)
		if err != nil {
			return nil, fmt.Errorf("invalid sub format: %v", err)
		}

		rolesStr, ok := claims["roles"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid roles format")
		}
		roles := strings.Split(rolesStr, ",")

		payload := &JwtPayload{
			Sub:      subID,
			Iat:      int64(claims["iat"].(float64)),
			Exp:      int64(claims["exp"].(float64)),
			Username: claims["username"].(string),
			Roles:    roles,
		}
		return payload, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}
