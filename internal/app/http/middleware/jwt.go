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
	"social/internal/usecase"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"social/internal/entity"
)

// ReleaseMode is the mode for release, such as dev or release.
type ReleaseMode string

const (
	// ReleaseModeProd is the prod mode.
	ReleaseModeProd ReleaseMode = "prod"
	// ReleaseModeDev is the dev mode.
	ReleaseModeDev ReleaseMode = "dev"
)

const (
	issuer                  = "bytebase"
	accessTokenAudienceFmt  = "bb.user.access.%s"
	refreshTokenAudienceFmt = "bb.user.refresh.%s"

	// Cookie section.
	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"

	// Signing key section. For now, this is only used for signing, not for verifying since we only
	// have 1 version. But it will be used to maintain backward compatibility if we change the signing mechanism.
	keyID = "v1"

	// Expiration section.
	refreshThresholdDuration = 1 * time.Hour
	accessTokenDuration      = 24 * time.Hour
	refreshTokenDuration     = 7 * 24 * time.Hour
	// Make cookie expire slightly earlier than the jwt expiration. Client would be logged out if the user
	// cookie expires, thus the client would always logout first before attempting to make a request with the expired jwt.
	// Suppose we have a valid refresh token, we will refresh the token in 2 cases:
	// 1. The access token is about to expire in <<refreshThresholdDuration>>
	// 2. The access token has already expired, we refresh the token so that the ongoing request can pass through.
	cookieExpDuration = refreshTokenDuration - 1*time.Minute

	// Context section
	// The key name used to store principal id in the context
	// principal id is extracted from the jwt token subject field.
	principalIDContextKey = "principal-id"
)

// Claims creates a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like name.
type Claims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func getPrincipalIDContextKey() string {
	return principalIDContextKey
}

// GenerateTokensAndSetCookies generates jwt token and saves it to the http-only cookie.
func GenerateTokensAndSetCookies(c *gin.Context, user *entity.User, mode ReleaseMode, secret string) error {
	accessToken, err := generateAccessToken(user, mode, secret)
	if err != nil {
		return fmt.Errorf("failed to generate access token: %w", err)
	}

	cookieExp := time.Now().Add(cookieExpDuration)
	setTokenCookie(c, accessTokenCookieName, accessToken, cookieExp.Second())
	setUserCookie(c, user, cookieExp.Second())

	// We generate here a new refresh token and saving it to the cookie.
	refreshToken, err := generateRefreshToken(user, mode, secret)
	if err != nil {
		return fmt.Errorf("failed to generate refresh token: %w", err)
	}
	setTokenCookie(c, refreshTokenCookieName, refreshToken, cookieExp.Second())

	return nil
}

func generateAccessToken(user *entity.User, mode ReleaseMode, secret string) (string, error) {
	expirationTime := time.Now().Add(accessTokenDuration)
	return generateToken(user, fmt.Sprintf(accessTokenAudienceFmt, mode), expirationTime, []byte(secret))
}

func generateRefreshToken(user *entity.User, mode ReleaseMode, secret string) (string, error) {
	expirationTime := time.Now().Add(refreshTokenDuration)
	return generateToken(user, fmt.Sprintf(refreshTokenAudienceFmt, mode), expirationTime, []byte(secret))
}

// Pay attention to this function. It holds the main JWT token generation logic.
func generateToken(user *entity.User, aud string, expirationTime time.Time, secret []byte) (string, error) {
	// Create the JWT claims, which includes the username and expiry time.
	claims := &Claims{
		Name: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience: jwt.ClaimStrings{aud},
			// In JWT, the expiry time is expressed as unix milliseconds.
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			Subject:   strconv.Itoa(user.ID),
		},
	}

	// Declare the token with the HS256 algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = keyID

	// Create the JWT string.
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Here we are creating a new cookie, which will store the valid JWT token.
func setTokenCookie(c *gin.Context, name, token string, expiration int) {
	c.SetCookie(name, token, expiration, "/", "localhost", false, true)
}

func removeTokenCookie(c *gin.Context, name string) {
	c.SetCookie(name, "", -1, "/", "localhost", false, true)
}

// Purpose of this cookie is to store the user's id.
func setUserCookie(c *gin.Context, user *entity.User, expiration int) {
	c.SetCookie("user", strconv.Itoa(user.ID), expiration, "/", "localhost", false, true)
}

func removeUserCookie(c *gin.Context) {
	c.SetCookie("user", "", -1, "/", "localhost", false, true)
}

// JWTMiddleware validates the access token.
// If the access token is about to expire or has expired and the request has a valid refresh token, it
// will try to generate new access token and refresh token.
func JWTMiddleware(userRepo *usecase.UserUseCase, mode ReleaseMode, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skips auth, actuator, plan
		if HasPrefixes(c.Request.RequestURI, "/api/auth", "/api/actuator", "/api/plan") {
			c.Next()
			return
		}

		//method := c.Request().Method
		//// Skip GET /subscription request
		//if HasPrefixes(c.Path(), "/api/subscription") && method == "GET" {
		//	return next(c)
		//}

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

		if !audienceContains(claims.Audience, fmt.Sprintf(accessTokenAudienceFmt, mode)) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 400,
				"msg": fmt.Sprintf("Invalid access token, audience mismatch, got %q, expected %q. you may send request to the wrong environment",
					claims.Audience,
					fmt.Sprintf(accessTokenAudienceFmt, mode),
				),
			})
			c.Abort()
			return
		}

		generateToken := time.Until(claims.ExpiresAt.Time) < refreshThresholdDuration
		if err != nil {
			var ve *jwt.ValidationError
			if errors.As(err, &ve) {
				// If expiration error is the only error, we will clear the err
				// and generate new access token and refresh token
				if ve.Errors == jwt.ValidationErrorExpired {
					err = nil
					generateToken = true
				}
			}
		}

		// We either have a valid access token or we will attempt to generate new access token and refresh token
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

			// Even if there is no error, we still need to make sure the user still exists.
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

					// Parses token and checks if it's valid.
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

					if !audienceContains(refreshTokenClaims.Audience, fmt.Sprintf(refreshTokenAudienceFmt, mode)) {
						return &generateTokenErr{
							code: http.StatusUnauthorized,
							msg: fmt.Sprintf("Invalid refresh token, audience mismatch, got %q, expected %q. you may send request to the wrong environment",
								refreshTokenClaims.Audience,
								fmt.Sprintf(refreshTokenAudienceFmt, mode),
							),
						}
					}

					// If we have a valid refresh token, we will generate new access token and refresh token
					if refreshToken != nil && refreshToken.Valid {
						if err := GenerateTokensAndSetCookies(c, user, mode, secret); err != nil {
							return &generateTokenErr{
								code: http.StatusInternalServerError,
								msg:  fmt.Sprintf("Server error to refresh expired token. User Id %d", uid),
							}
						}
					}

					return nil
				}

				// It may happen that we still have a valid access token, but we encounter issue when trying to generate new token
				// In such case, we won't return the error.
				if err := generateTokenFunc(); err != nil && !accessToken.Valid {
					c.JSON(http.StatusUnauthorized, gin.H{
						"code": err.code,
						"msg":  err.msg,
					})
					c.Abort()
					return
				}
			}

			// Stores principalID into context.
			c.Set(getPrincipalIDContextKey(), uid)
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

// FindString returns the search index of sorted strings.
func FindString(stringList []string, search string) int {
	sort.Strings(stringList)
	i := sort.SearchStrings(stringList, search)
	if i == len(stringList) {
		return -1
	}
	return i
}

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandomString returns a random string with length n.
func RandomString(n int) (string, error) {
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		// The reason for using crypto/rand instead of math/rand is that
		// the former relies on hardware to generate random numbers and
		// thus has a stronger source of random numbers.
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

// HasPrefixes returns true if the string s has any of the given prefixes.
func HasPrefixes(src string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(src, prefix) {
			return true
		}
	}
	return false
}

// GetPostgresDataDir returns the postgres data directory of Bytebase.
func GetPostgresDataDir(dataDir string) string {
	return path.Join(dataDir, "pgdata")
}

// GetPostgresSocketDir returns the postgres socket directory of Bytebase.
func GetPostgresSocketDir() string {
	return "/tmp"
}

// GetResourceDir returns the resource directory of Bytebase.
func GetResourceDir(dataDir string) string {
	return path.Join(dataDir, "resources")
}

// DefaultMigrationVersion returns the default migration version string.
// Use the current time in second to guarantee uniqueness in a monotonic increasing way.
// We cannot add task ID because tenant mode databases should use the same migration version string when applying a schema update.
func DefaultMigrationVersion() string {
	return time.Now().Format("20060102150405")
}

// ParseTemplateTokens parses the template and returns template tokens and their delimiters.
// For example, if the template is "{{DB_NAME}}_hello_{{LOCATION}}", then the tokens will be ["{{DB_NAME}}", "{{LOCATION}}"],
// and the delimiters will be ["_hello_"].
// The caller will usually replace the tokens with a normal string, or a regexp. In the latter case, it will be a problem
// if there are special regexp characters like "$" in the delimiters. The caller should escape the delimiters in such cases.
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

// GetFileSizeSum calculates the sum of file sizes for file names in the list.
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
