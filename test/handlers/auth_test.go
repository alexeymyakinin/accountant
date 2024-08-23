package handlers

import (
	"accountant/internal/api"
	"accountant/internal/config"
	"accountant/internal/dependencies"
	"accountant/test"
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"testing"
)

type authSuite struct {
	suite.Suite
	app       *fiber.App
	container *dependencies.Container
}

func (s *authSuite) SetupSuite() {
	s.Require().NoError(config.LoadConfig())

	p, err := test.SetupDatabaseAndPool(config.Get().DbUrl, "test")
	s.Require().NoError(err)

	c := dependencies.New(p, config.Get())
	s.app = api.New(c)
	s.container = c
}

func (s *authSuite) TearDownSuite() {
	s.container.Pool.Close()
	s.Require().NoError(s.app.Shutdown())
	s.Require().NoError(test.TeardownDatabase(config.Get().DbUrl, "test"))
}

func (s *authSuite) TearDownTest() {
	s.Require().NoError(test.DeleteFromTables(s.container.Pool, "users"))
}

func (s *authSuite) TestOk() {
	req, err := test.NewFormRequest(http.MethodPost, "/register", map[string]any{
		"email":    "test@test.com",
		"password": "12345678",
	})
	s.Require().NoError(err)

	resp, err := s.app.Test(req)
	s.Require().NoError(err)

	data, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	body, err := test.JsonToMap(data)
	s.Require().NoError(err, string(data))

	s.Require().Equal(http.StatusOK, resp.StatusCode, string(data))
	s.Require().NotEmpty(body["token"], body)
}

func (s *authSuite) TestEmailEmpty() {
	req, err := test.NewFormRequest(http.MethodPost, "/register", map[string]any{
		"email":    "",
		"password": "12345678",
	})
	s.Require().NoError(err)

	resp, err := s.app.Test(req)
	s.Require().NoError(err)

	data, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	body, err := test.JsonToMap(data)
	s.Require().NoError(err, string(data))

	s.Require().Equal(http.StatusBadRequest, resp.StatusCode, string(data))
	s.Require().Equal("Key: 'registerRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag", body["detail"])
}

func (s *authSuite) TestEmailIncorrect() {
	req, err := test.NewFormRequest(http.MethodPost, "/register", map[string]any{
		"email":    "test",
		"password": "12345678",
	})
	s.Require().NoError(err)

	resp, err := s.app.Test(req)
	s.Require().NoError(err)

	data, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	body, err := test.JsonToMap(data)
	s.Require().NoError(err, string(data))

	s.Require().Equal(http.StatusBadRequest, resp.StatusCode, string(data))
	s.Require().Equal("Key: 'registerRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag", body["detail"])
}

func (s *authSuite) TestInvalidPassword() {
	req, err := test.NewFormRequest(http.MethodPost, "/register", map[string]any{
		"email":    "test@test.com",
		"password": "",
	})
	s.Require().NoError(err)

	resp, err := s.app.Test(req)
	s.Require().NoError(err)

	data, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	body, err := test.JsonToMap(data)
	s.Require().NoError(err, string(data))

	s.Require().Equal(http.StatusBadRequest, resp.StatusCode, string(data))
	s.Require().Equal("Key: 'registerRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag", body["detail"])
}

func (s *authSuite) TestEmailExists() {
	_, err := s.container.Pool.Exec(
		context.Background(),
		"INSERT INTO users(email, password) VALUES ($1, $2)",
		"test@test.com",
		"12345678",
	)
	s.Require().NoError(err)

	req, err := test.NewFormRequest(http.MethodPost, "/register", map[string]any{
		"email":    "test@test.com",
		"password": "12345678",
	})
	s.Require().NoError(err)

	resp, err := s.app.Test(req)
	s.Require().NoError(err)

	data, err := io.ReadAll(resp.Body)
	s.Require().NoError(err)

	body, err := test.JsonToMap(data)
	s.Require().NoError(err, string(data))

	s.Require().Equal(http.StatusBadRequest, resp.StatusCode, string(data))
	s.Require().Equal("email already taken", body["detail"])
}

func TestAuth(t *testing.T) {
	suite.Run(t, new(authSuite))
}
