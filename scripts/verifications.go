package scripts

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/ali-shokoohi/micro-gopia/config"
	"github.com/ali-shokoohi/micro-gopia/internal/domain/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// CheckPassword checks passwords must be bigger than 8 character with uppercase, lowercase, digits and symbols
func CheckPassword(s string) bool {
	letters := 0
	eightOrMore, number, upper, special := false, false, false, false
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
		}
	}
	eightOrMore = letters >= 8
	if eightOrMore && number && upper && special {
		return true
	}
	return false
}

func ValidateToken(c *gin.Context) error {
	_, _, err := GetToken(c)
	return err

}

func GetToken(c *gin.Context) (*jwt.Token, *dto.Claims, error) {
	tokenString := getTokenFromRequest(c)
	tk := &dto.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, tk, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Confs.Service.Token.Password), nil
	})
	return token, tk, err
}

func getTokenFromRequest(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")

	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}

func CurrentTokenClaim(c *gin.Context) (*dto.Claims, error) {
	err := ValidateToken(c)
	if err != nil {
		return nil, err
	}
	_, tk, err := GetToken(c)
	if err != nil {
		return nil, err
	}
	return tk, nil
}
