package services

import (
	"accountant/internal/domain"
	"accountant/internal/infra/repositories"
	"accountant/internal/infra/uow"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	UserNotFoundError       = errors.New("user not found")
	EmailAlreadyExistsError = errors.New("email already taken")
	IncorrectPasswordError  = errors.New("incorrect password")
)

type AuthService struct {
	Uow              *uow.UnitOfWork
	Redis            *redis.Client
	JwtDuration      time.Duration
	JwtSigningKey    []byte
	JwtSigningMethod jwt.SigningMethod
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := uow.Execute(ctx, s.Uow, func(c *uow.Context) (domain.User, error) {
		return c.UserRepository.GetByEmail(ctx, email)
	})
	if err != nil {
		switch {
		case errors.Is(err, repositories.NotFoundError):
			return "", UserNotFoundError
		default:
			return "", err
		}
	}

	if !s.isCorrectPassword(user.HashedPassword, password) {
		return "", IncorrectPasswordError
	}

	return s.encodeToken(s.createToken(user))
}

func (s *AuthService) Register(ctx context.Context, email string, password string) (string, error) {
	opts := pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	}

	user, err := uow.ExecuteTx(ctx, s.Uow, func(c *uow.Context) (domain.User, error) {
		exist, err := uow.Execute(ctx, s.Uow, func(c *uow.Context) (bool, error) {
			return c.UserRepository.IsEmailExists(ctx, email)
		})
		if err != nil {
			return domain.User{}, err
		}

		if exist {
			return domain.User{}, EmailAlreadyExistsError
		}

		hashedPassword, err := s.hashedPassword(password)
		if err != nil {
			return domain.User{}, err
		}

		id, err := uow.Execute(ctx, s.Uow, func(c *uow.Context) (int, error) {
			return c.UserRepository.Insert(ctx, email, hashedPassword)
		})
		if err != nil {
			return domain.User{}, err
		}

		return domain.User{Id: id, Email: email}, nil
	}, opts)
	if err != nil {
		return "", err
	}

	return s.encodeToken(s.createToken(user))
}

func (s *AuthService) CreateToken(user domain.User) (string, error) {
	return s.encodeToken(s.createToken(user))
}

func (s *AuthService) isCorrectPassword(hashedPassword string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func (s *AuthService) hashedPassword(password string) (string, error) {
	r, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(r), nil
}

func (s *AuthService) createToken(user domain.User) *jwt.Token {
	return jwt.NewWithClaims(s.JwtSigningMethod, domain.JWTClaims{
		UserId:    user.Id,
		UserEmail: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.JwtDuration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})
}

func (s *AuthService) encodeToken(token *jwt.Token) (string, error) {
	r, err := token.SignedString(s.JwtSigningKey)
	if err != nil {
		return "", fmt.Errorf("error during encoding token: %w", err)
	}

	return r, nil
}

func (s *AuthService) decodeToken(tokenString string) (domain.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.JwtSigningKey, nil
	})
	if err != nil {
		return domain.JWTClaims{}, err
	}

	return token.Claims.(domain.JWTClaims), nil
}
