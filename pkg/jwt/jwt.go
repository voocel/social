package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"social/ent"
)

const (
	issuer                  = "social"
	accessTokenAudienceFmt  = "social.user.access"
	refreshTokenAudienceFmt = "social.user.refresh"

	keyID     = "v1"
	jwtSecret = "social-key"

	refreshThresholdDuration = 1 * time.Hour
	accessTokenDuration      = 24 * time.Hour
	refreshTokenDuration     = 7 * 24 * time.Hour

	socialContextKey = "user"
)

type Claims struct {
	User *ent.User `json:"user"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 jwt token
func GenerateToken(user *ent.User) (accessToken, refreshToken string, err error) {
	accessToken, err = generateAccessToken(user)
	if err != nil {
		err = fmt.Errorf("failed to generate access token: %w", err)
		return
	}

	refreshToken, err = generateRefreshToken(user)
	if err != nil {
		err = fmt.Errorf("failed to generate refresh token: %w", err)
		return
	}
	return
}

func generateAccessToken(user *ent.User) (string, error) {
	expirationTime := time.Now().Add(accessTokenDuration)
	return generateToken(user, accessTokenAudienceFmt, expirationTime)
}

func generateRefreshToken(user *ent.User) (string, error) {
	expirationTime := time.Now().Add(refreshTokenDuration)
	return generateToken(user, refreshTokenAudienceFmt, expirationTime)
}

func generateToken(user *ent.User, aud string, expirationTime time.Time) (string, error) {
	claims := &Claims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{aud},
			ExpiresAt: jwt.NewNumericDate(expirationTime), // unix milliseconds
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			Subject:   strconv.Itoa(int(user.ID)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = keyID
	return token.SignedString([]byte(jwtSecret))
}

func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, err
	} else {
		return nil, err
	}
}
