package dependencies

import (
	"accountant/internal/config"
	"accountant/internal/infra/uow"
	"accountant/internal/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	Pool   *pgxpool.Pool
	Config *config.Config
}

func New(pool *pgxpool.Pool, conf *config.Config) *Container {
	return &Container{Pool: pool, Config: conf}
}

func (c *Container) Uow() *uow.UnitOfWork {
	return uow.New(c.Pool)
}

func (c *Container) AuthService() *services.AuthService {
	return &services.AuthService{
		Uow:              c.Uow(),
		JwtDuration:      c.Config.JwtDuration,
		JwtSigningKey:    []byte(c.Config.JwtSigningKey),
		JwtSigningMethod: jwt.SigningMethodHS256,
	}
}

func (c *Container) AccountService() *services.AccountService {
	return &services.AccountService{
		Uow: c.Uow(),
	}
}
