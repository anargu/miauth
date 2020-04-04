package miauth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)


const defaultExpirationOffsetTime = 5 * 60 * 60
func expirationOffset(offset *int64) int64 {
	if offset != nil {
		return time.Now().Unix() + *offset
	} else {
		return time.Now().Unix() + defaultExpirationOffsetTime
	}
}

type CustomClaims struct {
	UserID string `json:"userId"`
	UserEmail string `json:"user_email"`
	jwt.StandardClaims
}
func tokenize(userId string, userEmail string, secret string, expires *int64) (string,error) {
	var claims CustomClaims
	if expires == nil {
		claims = CustomClaims{
			UserID: userId,
			UserEmail: userEmail,
		}
	} else {
		claims = CustomClaims{
			UserID: userId,
			UserEmail: userEmail,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: *expires,
			},
		}
	}
	preToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := preToken.SignedString([]byte(secret))
	return token, err
}

func TokenizeAccessToken(userId string, userEmail string) (string, string, error) {
	expIn, err := strconv.ParseInt(Config.AccessToken.ExpiresIn, 10, 64)
	if err != nil {
		return "", "", err
	}
	exp := expirationOffset(&expIn)
	token, err := tokenize(userId, userEmail, Config.AccessToken.Secret, &exp)
	return token, strconv.FormatInt(exp, 10),err
}

func TokenizeRefreshToken(userId string, userEmail string) (string, error) {
	return tokenize(userId, userEmail, Config.RefreshToken.Secret, nil)
}

func TokenizeResetEmailToken(userId string, userEmail string) (string, error)  {
	expIn, err := strconv.ParseInt(Config.ResetPassword.ExpiresIn, 10, 64)
	if err != nil {
		return "", err
	}
	exp := expirationOffset(&expIn)
	return tokenize(userId, userEmail, Config.ResetPassword.Secret, &exp)
}


func verify(token string, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
}

func VerifyAccessToken(token string) (bool, error) {
	tokenResult, err := verify(token, Config.AccessToken.Secret)
	if err != nil {
		return false, err
	}
	return tokenResult.Valid, nil
}

func VerifyRefreshToken(token string) (bool, error) {
	tokenResult, err := verify(token, Config.RefreshToken.Secret)
	if err != nil {
		return false, err
	}
	return tokenResult.Valid, nil
}

func VerifyResetEmailToken(token string) (*string, error) {
	tokenResult, err := verify(token, Config.ResetPassword.Secret)
	if err != nil {
		return nil, err
	}

	if claims, ok := tokenResult.Claims.(CustomClaims); ok && tokenResult.Valid {
		return &claims.UserID, nil
	} else {
		return nil, errors.New("invalid JWT Token")
	}
}

func ComparePassword(planPassword string, hashedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(planPassword)); err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) (*string, error) {
	salt, err := strconv.Atoi(Config.BCrypt.Salt)
	if err != nil {
		return nil, err
	}
	hashInBytes, err := bcrypt.GenerateFromPassword([]byte(password), salt)
	if err != nil {
		return nil, err
	}
	hash := string(hashInBytes)
	return &hash, nil
}