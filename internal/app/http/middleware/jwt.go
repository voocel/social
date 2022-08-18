package middleware

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"social/internal/entity"
	"social/internal/usecase"
)

const (
	issuer                  = "social"
	accessTokenAudienceFmt  = "social.user.access"
	refreshTokenAudienceFmt = "social.user.refresh"

	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"

	keyID = "v1"

	refreshThresholdDuration = 1 * time.Hour
	accessTokenDuration      = 24 * time.Hour
	refreshTokenDuration     = 7 * 24 * time.Hour

	// cookie的过期时间略早于jwt过期时间, 使得如果用户cookie过期，则会在使用过期的jwt前先注销
	cookieExpDuration = refreshTokenDuration - 1*time.Minute

	socialIDContextKey = "social-id"
)

type Claims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func getSocialIDContextKey() string {
	return socialIDContextKey
}

// GenerateTokensAndSetCookies 生成 jwt token 并且保存到 http-only cookie
func GenerateTokensAndSetCookies(c *gin.Context, user *entity.User, secret string) error {
	accessToken, err := generateAccessToken(user, secret)
	if err != nil {
		return fmt.Errorf("failed to generate access token: %w", err)
	}

	cookieExp := time.Now().Add(cookieExpDuration)
	setTokenCookie(c, accessTokenCookieName, accessToken, cookieExp.Second())
	setUserCookie(c, user, cookieExp.Second())

	// 生成refreshToken并保存到cookie中
	refreshToken, err := generateRefreshToken(user, secret)
	if err != nil {
		return fmt.Errorf("failed to generate refresh token: %w", err)
	}
	setTokenCookie(c, refreshTokenCookieName, refreshToken, cookieExp.Second())

	return nil
}

func generateAccessToken(user *entity.User, secret string) (string, error) {
	expirationTime := time.Now().Add(accessTokenDuration)
	return generateToken(user, fmt.Sprintf(accessTokenAudienceFmt), expirationTime, []byte(secret))
}

func generateRefreshToken(user *entity.User, secret string) (string, error) {
	expirationTime := time.Now().Add(refreshTokenDuration)
	return generateToken(user, refreshTokenAudienceFmt, expirationTime, []byte(secret))
}

func generateToken(user *entity.User, aud string, expirationTime time.Time, secret []byte) (string, error) {
	claims := &Claims{
		Name: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{aud},
			ExpiresAt: jwt.NewNumericDate(expirationTime), // unix milliseconds
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			Subject:   strconv.Itoa(user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = keyID
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 创建一个新的cookie，并存储有效的token
func setTokenCookie(c *gin.Context, name, token string, expiration int) {
	c.SetCookie(name, token, expiration, "/", "localhost", false, true)
}

// 清除cookie
func removeTokenCookie(c *gin.Context, name string) {
	c.SetCookie(name, "", -1, "/", "localhost", false, true)
}

// 存储用户id
func setUserCookie(c *gin.Context, user *entity.User, expiration int) {
	c.SetCookie("user", strconv.Itoa(user.ID), expiration, "/", "localhost", false, true)
}

func removeUserCookie(c *gin.Context) {
	c.SetCookie("user", "", -1, "/", "localhost", false, true)
}

// JWTMiddleware 验证 access token
func JWTMiddleware(userRepo *usecase.UserUseCase, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过 auth, register
		if HasPrefixes(c.Request.RequestURI, "/api/auth", "/api/register", "/api/plan") {
			c.Next()
			return
		}

		cookie, err := c.Cookie(accessTokenCookieName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 400,
				"msg":  "Missing access token",
			})
			c.Abort()
			return
		}

		claims := &Claims{}
		accessToken, err := jwt.ParseWithClaims(cookie, claims, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Name {
				return nil, fmt.Errorf("unexpected access token signing method=%v, expect %v", t.Header["alg"], jwt.SigningMethodHS256)
			}
			if kid, ok := t.Header["kid"].(string); ok {
				if kid == "v1" {
					return []byte(secret), nil
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

		generateToken := time.Until(claims.ExpiresAt.Time) < refreshThresholdDuration
		if err != nil {
			var ve *jwt.ValidationError
			if errors.As(err, &ve) {
				// token过期则清除错误并刷新token
				if ve.Errors == jwt.ValidationErrorExpired {
					err = nil
					generateToken = true
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
			user, err := userRepo.GetUserById(ctx, uid)
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

			if generateToken {
				generateTokenFunc := func() *generateTokenErr {
					rc, err := c.Cookie(refreshTokenCookieName)

					if err != nil {
						return &generateTokenErr{
							code: http.StatusUnauthorized,
							msg:  "Failed to generate access token. Missing refresh token.",
						}
					}

					// 解析token并检查其是否有效.
					refreshTokenClaims := &Claims{}
					refreshToken, err := jwt.ParseWithClaims(rc, refreshTokenClaims, func(t *jwt.Token) (interface{}, error) {
						if t.Method.Alg() != jwt.SigningMethodHS256.Name {
							return nil, fmt.Errorf("unexpected refresh token signing method=%v, expected %v", t.Header["alg"], jwt.SigningMethodHS256)
						}

						if kid, ok := t.Header["kid"].(string); ok {
							if kid == "v1" {
								return []byte(secret), nil
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
						if err := GenerateTokensAndSetCookies(c, user, secret); err != nil {
							return &generateTokenErr{
								code: http.StatusInternalServerError,
								msg:  fmt.Sprintf("Server error to refresh expired token. User Id %d", uid),
							}
						}
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

			// 将socialID存储到上下文
			c.Set(getSocialIDContextKey(), uid)
			c.Next()
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

func GetFileSizeSum(fileNameList []string) (int64, error) {
	var sum int64
	for _, fileName := range fileNameList {
		stat, err := os.Stat(fileName)
		if err != nil {
			return 0, err
		}
		sum += stat.Size()
	}
	return sum, nil
}
