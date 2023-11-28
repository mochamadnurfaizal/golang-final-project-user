package helpers

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = "zollav"

var invalidTokens = make(map[string]struct{})
var mu sync.Mutex

func AddInvalidToken(token string) {
	mu.Lock()
	defer mu.Unlock()
	invalidTokens[token] = struct{}{}
}

func IsTokenInvalid(token string) bool {
	mu.Lock()
	defer mu.Unlock()
	_, invalid := invalidTokens[token]
	return invalid
}

type Claims struct {
	ID       uint
	Username string
	Fullname string
	Roles    string
	jwt.StandardClaims
}

func GenerateToken(id uint, username string, fullname string, roles string) string {
	secretKey = os.Getenv("SECRET_KEY")
	claims := Claims{
		ID:       id,
		Username: username,
		Fullname: fullname,
		Roles:    roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, _ := parseToken.SignedString([]byte(secretKey))
	return signedToken
}

func VerifyToken(tokenString string) (*Claims, error) {
	secretKey = os.Getenv("SECRET_KEY")

	errResponse := errors.New("Token Can't be Used")
	bearer := strings.HasPrefix(tokenString, "Bearer")

	if !bearer {
		return nil, errResponse
	}

	stringToken := strings.Split(tokenString, " ")[1]
	mu.Lock()
	defer mu.Unlock()
	if _, invalid := invalidTokens[stringToken]; invalid {
		return nil, fmt.Errorf("token is invalid")
	}

	token, err := jwt.ParseWithClaims(stringToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func HashPass(p string) string {
	salt := 8
	password := []byte(p)
	hash, _ := bcrypt.GenerateFromPassword(password, salt)

	return string(hash)
}

func ComparePass(h, p []byte) bool {
	hash, pass := []byte(h), []byte(p)

	err := bcrypt.CompareHashAndPassword(hash, pass)

	return err == nil
}

func GetContentType(ctx echo.Context) string {
	return ctx.Request().Header.Get("Content-Type")
}

func ConvertMapClaimsToMap(claims jwt.MapClaims) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range claims {
		result[key] = value
	}
	return result
}
