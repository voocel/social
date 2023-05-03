package middleware

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"social/ent"
	"social/internal/usecase"
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

	// cookie的过期时间略早于jwt过期时间, 使得如果用户cookie过期，则会在使用过期的jwt前先注销
	cookieExpDuration = refreshTokenDuration - 1*time.Minute

	socialContextKey = "jwt-user"
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

// JWTMiddleware 验证 access token
func JWTMiddleware(userRepo *usecase.UserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过 auth, register
		if HasPrefixes(c.Request.RequestURI, "/api/auth", "/v1/user/register", "/v1/user/login", "/swagger") {
			c.Next()
			return
		}

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 400,
				"msg":  "Missing access token",
			})
			c.Abort()
			return
		}

		claims := &Claims{}
		accessToken, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Name {
				return nil, fmt.Errorf("unexpected access token signing method=%v, expect %v", t.Header["alg"], jwt.SigningMethodHS256)
			}
			if kid, ok := t.Header["kid"].(string); ok {
				if kid == "v1" {
					return []byte(jwtSecret), nil
				}
			}
			return nil, fmt.Errorf("unexpected access token kid=%v", t.Header["kid"])
		})

		if !audienceContains(claims.Audience, accessTokenAudienceFmt) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 400,
				"msg": fmt.Sprintf("Invalid access token, audience mismatch, got %q, expected %q. you may send request to the wrong environment",
					claims.Audience,
					accessTokenAudienceFmt,
				),
			})
			c.Abort()
			return
		}

		isGenToken := time.Until(claims.ExpiresAt.Time) < refreshThresholdDuration
		if err != nil {
			var ve *jwt.ValidationError
			if errors.As(err, &ve) {
				// token过期则清除错误并刷新token
				if ve.Errors == jwt.ValidationErrorExpired {
					err = nil
					isGenToken = true
				}
			}
		}

		// 要么拥有有效的访问令牌，要么尝试生成新的访问令牌和刷新令牌
		if err == nil {
			ctx := c.Request.Context()
			uid, err := strconv.Atoi(claims.Subject)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": 400,
					"msg":  "Malformed ID in the token",
				})
				c.Abort()
				return
			}

			// 即使没有错误，仍然需要确保用户仍然存在。
			user, err := userRepo.GetUserById(ctx, int64(uid))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 401,
					"msg":  fmt.Sprintf("Server error to find user ID: %d", uid),
				})
				c.Abort()
				return
			}
			if user == nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code": 400,
					"msg":  fmt.Sprintf("Failed to find user ID: %d", uid),
				})
				c.Abort()
				return
			}

			if isGenToken {
				generateTokenFunc := func() *generateTokenErr {
					refreshTokenString := c.GetHeader("refresh-token")
					if refreshTokenString == "" {
						return &generateTokenErr{
							code: http.StatusUnauthorized,
							msg:  "Failed to generate access token. Missing refresh token.",
						}
					}

					// 解析token并检查其是否有效.
					refreshTokenClaims := &Claims{}
					refreshToken, err := jwt.ParseWithClaims(refreshTokenString, refreshTokenClaims, func(t *jwt.Token) (interface{}, error) {
						if t.Method.Alg() != jwt.SigningMethodHS256.Name {
							return nil, fmt.Errorf("unexpected refresh token signing method=%v, expected %v", t.Header["alg"], jwt.SigningMethodHS256)
						}

						if kid, ok := t.Header["kid"].(string); ok {
							if kid == "v1" {
								return []byte(jwtSecret), nil
							}
						}
						return nil, fmt.Errorf("unexpected refresh token kid=%v", t.Header["kid"])
					})
					if err != nil {
						if err == jwt.ErrSignatureInvalid {
							return &generateTokenErr{
								code: http.StatusUnauthorized,
								msg:  "Failed to generate access token. Invalid refresh token signature.",
							}
						}
						return &generateTokenErr{
							code: http.StatusInternalServerError,
							msg:  fmt.Sprintf("Server error to refresh expired token. User Id %d", uid),
						}
					}

					if !audienceContains(refreshTokenClaims.Audience, refreshTokenAudienceFmt) {
						return &generateTokenErr{
							code: http.StatusUnauthorized,
							msg: fmt.Sprintf("Invalid refresh token, audience mismatch, got %q, expected %q. you may send request to the wrong environment",
								refreshTokenClaims.Audience,
								refreshTokenAudienceFmt,
							),
						}
					}

					// 如果有一个有效的refresh token，将生成新的access token和refresh token
					if refreshToken != nil && refreshToken.Valid {
						at, _, err := GenerateToken(user)
						if err != nil {
							return &generateTokenErr{
								code: http.StatusInternalServerError,
								msg:  fmt.Sprintf("Server error to refresh expired token. User Id %d", uid),
							}
						}
						c.JSON(http.StatusOK, gin.H{
							"code": 0,
							"msg":  "refresh token success",
							"data": map[string]string{
								"access_token": at,
							},
						})
						c.Abort()
					}

					return nil
				}

				// 这可能仍然有一个有效的access token，但在尝试生成新令牌时会遇到问题，这种情况不反回错误
				if err := generateTokenFunc(); err != nil && !accessToken.Valid {
					c.JSON(http.StatusUnauthorized, gin.H{
						"code": err.code,
						"msg":  err.msg,
					})
					c.Abort()
					return
				}
			}

			// 将用户信息存储到上下文
			c.Set(socialContextKey, user)
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "Invalid or expired access token",
		})
		c.Abort()
		return
	}
}

type generateTokenErr struct {
	code int
	msg  string
}

func audienceContains(audience jwt.ClaimStrings, token string) bool {
	for _, v := range audience {
		if v == token {
			return true
		}
	}
	return false
}

func FindString(stringList []string, search string) int {
	sort.Strings(stringList)
	i := sort.SearchStrings(stringList, search)
	if i == len(stringList) {
		return -1
	}
	return i
}

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(n int) (string, error) {
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		randNum, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		if _, err := sb.WriteRune(letters[randNum.Uint64()]); err != nil {
			return "", err
		}
	}
	return sb.String(), nil
}

func HasPrefixes(src string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(src, prefix) {
			return true
		}
	}
	return false
}

func GetPostgresDataDir(dataDir string) string {
	return path.Join(dataDir, "pgdata")
}

func GetPostgresSocketDir() string {
	return "/tmp"
}

func GetResourceDir(dataDir string) string {
	return path.Join(dataDir, "resources")
}

func DefaultMigrationVersion() string {
	return time.Now().Format("20060102150405")
}

// ParseTemplateTokens 解析模板并返回模板标记及其分隔符
func ParseTemplateTokens(template string) ([]string, []string) {
	r := regexp.MustCompile(`{{[^{}]+}}`)
	tokens := r.FindAllString(template, -1)
	if len(tokens) > 0 {
		split := r.Split(template, -1)
		var delimiters []string
		for _, s := range split {
			if s != "" {
				delimiters = append(delimiters, s)
			}
		}
		return tokens, delimiters
	}
	return nil, nil
}
