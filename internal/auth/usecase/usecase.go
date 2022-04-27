package usecase

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strings"
	"time"

	"github.com/alexkopcak/gophermart/internal/auth"
	"github.com/alexkopcak/gophermart/internal/models"
	"github.com/dgrijalva/jwt-go/v4"
)

type AuthClaims struct {
	jwt.StandardClaims
	User *models.User
}

type AuthUseCase struct {
	userRepo       auth.UserRepository
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthUseCase(userRepo auth.UserRepository,
	hashSalt string,
	signingKey string,
	tokenTTL time.Duration) auth.UseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		hashSalt:       hashSalt,
		signingKey:     []byte(signingKey),
		expireDuration: tokenTTL,
	}
}

func generatePasswordHash(password string, hashSalt string) string {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(hashSalt))

	return fmt.Sprintf("%x", pwd.Sum(nil))
}

func (auc *AuthUseCase) SignUp(ctx context.Context, userName string, password string) error {
	user := &models.User{
		UserName: userName,
		Password: generatePasswordHash(password, auc.hashSalt),
	}
	return auc.userRepo.CreateUser(ctx, user)
}

func (auc *AuthUseCase) SignIn(ctx context.Context, userName string, password string) (string, error) {
	pwd := generatePasswordHash(password, auc.hashSalt)

	user, err := auc.userRepo.GetUser(ctx, userName)
	if err != nil {
		return "", err
	}

	if strings.Compare(pwd, user.Password) != 0 {
		return "", auth.ErrBadLoginPassword
	}

	claims := AuthClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(auc.expireDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(auc.signingKey)
}

func (auc *AuthUseCase) ParseToken(ctx context.Context, accessToken string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, auth.ErrInternalServer
			}
			return auc.signingKey, nil
		})

	if err != nil {
		return nil, auth.ErrInternalServer
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims.User, nil
	}

	return nil, auth.ErrInternalServer
}
