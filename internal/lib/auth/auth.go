package auth

import (
	"time"

	"github.com/leagueify/leagueify/internal/config"
	"github.com/leagueify/leagueify/internal/lib/error"
	"github.com/leagueify/leagueify/internal/lib/token"
	"github.com/leagueify/leagueify/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type leagueifyJWTClaims struct {
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

func CreateJWT(account model.Account, aud string, maxAge int) (string, error) {
	cfg := config.LoadConfig()

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, leagueifyJWTClaims{
			account.IsAdmin,
			jwt.RegisteredClaims{
				Subject:  token.ReturnSignedToken(account.ID),
				Audience: []string{aud},
				ExpiresAt: jwt.NewNumericDate(
					time.Now().Add(time.Duration(maxAge) * time.Minute),
				),
				NotBefore: jwt.NewNumericDate(
					time.Now().Add(-1 * time.Minute),
				),
				IssuedAt: jwt.NewNumericDate(
					time.Now(),
				),
			},
		},
	)

	signedToken, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func HashPassword(providedPassword *string) error {
	if err := passwordRequirements(providedPassword); err != nil {
		return err
	}
	password := []byte(*providedPassword)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 12)
	if err != nil {
		return err
	}
	*providedPassword = string(hashedPassword)
	return nil
}

func PasswordsMatch(password, accountPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(accountPassword), []byte(password),
	)
	return err == nil
}

func VerifyJWT(tokenString string) (*leagueifyJWTClaims, error) {
	cfg := config.LoadConfig()
	token, err := jwt.ParseWithClaims(
		tokenString, &leagueifyJWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, &errors.LeagueifyError{
					Message: "invalid auth token",
				}
			}
			return []byte(cfg.JWTSecret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*leagueifyJWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, &errors.LeagueifyError{Message: "invalid auth token"}
}

func passwordRequirements(providedPassword *string) error {
	if len(*providedPassword) < 8 || len(*providedPassword) > 64 {
		return &errors.LeagueifyError{
			Message: "password must be 8-64 characters long",
		}
	}
	return nil
}
