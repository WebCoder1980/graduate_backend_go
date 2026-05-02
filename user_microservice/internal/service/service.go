package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"user_microservice/internal/model"

	"github.com/Nerzal/gocloak/v13"
)

type KeycloakClient interface {
	Login(ctx context.Context, clientID, clientSecret, realm, username, password string) (*gocloak.JWT, error)
	RefreshToken(ctx context.Context, refreshToken, clientID, clientSecret, realm string) (*gocloak.JWT, error)
	CreateUser(ctx context.Context, token, realm string, user gocloak.User) (string, error)
	SetPassword(ctx context.Context, token, userID, realm, password string, temporary bool) error
	DeleteUser(ctx context.Context, token, realm, id string) error
	GetRealmRole(ctx context.Context, token, realm, roleName string) (*gocloak.Role, error)
	CreateRealmRole(ctx context.Context, token, realm string, role gocloak.Role) (string, error)
	AddRealmRoleToUser(ctx context.Context, token, realm, userID string, roles []gocloak.Role) error
	LoginAdmin(ctx context.Context, username, password, realm string) (*gocloak.JWT, error)
}

type Service struct {
	ctx      context.Context
	keycloak KeycloakClient
}

func NewService(ctx context.Context) (*Service, error) {
	res := Service{
		ctx:      ctx,
		keycloak: gocloak.NewClient(os.Getenv("keycloak_address")),
	}

	err := res.init()
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func NewServiceWithKeycloak(ctx context.Context, kc KeycloakClient) (*Service, error) {
	res := Service{
		ctx:      ctx,
		keycloak: kc,
	}

	err := res.init()
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (s *Service) init() error {
	adminToken, err := s.loginAdmin()
	if err != nil {
		return err
	}

	_, err = s.validateRole(adminToken, "user")
	if err != nil {
		return err
	}

	_, err = s.validateRole(adminToken, "admin")
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) validateRole(adminToken *gocloak.JWT, name string) (*gocloak.Role, error) {
	role := gocloak.Role{
		Name:        gocloak.StringP(name),
		Description: gocloak.StringP(fmt.Sprintf("${role_%s}", name)),
	}
	_, _ = s.keycloak.CreateRealmRole(s.ctx, adminToken.AccessToken, os.Getenv("keycloak_realm"), role)

	return s.keycloak.GetRealmRole(s.ctx, adminToken.AccessToken, os.Getenv("keycloak_realm"), name)
}

func (s *Service) UserLogin(body *model.UserLogin) (*gocloak.JWT, error) {
	jwt, err := s.keycloak.Login(s.ctx, os.Getenv("keycloak_client_id"), os.Getenv("keycloak_client_secret"), os.Getenv("keycloak_realm"), body.Username, body.Password)
	if err != nil {
		return nil, err
	}

	return jwt, nil
}

func (s *Service) UserRefreshToken(body *model.UserRefreshToken) (*gocloak.JWT, error) {
	jwt, err := s.keycloak.RefreshToken(s.ctx, body.RefreshToken, os.Getenv("keycloak_client_id"), os.Getenv("keycloak_client_secret"), os.Getenv("keycloak_realm"))
	if err != nil {
		return nil, err
	}

	return jwt, nil
}

func (s *Service) UserRegisterPost(body *model.UserRegisterRequest) error {
	adminToken, err := s.loginAdmin()
	if err != nil {
		return err
	}

	newUserModel := gocloak.User{
		Enabled:  gocloak.BoolP(true),
		Username: gocloak.StringP(body.Username),
	}

	userId, err := s.keycloak.CreateUser(s.ctx, adminToken.AccessToken, os.Getenv("keycloak_realm"), newUserModel)
	if err != nil {
		return err
	}

	err = s.keycloak.SetPassword(s.ctx, adminToken.AccessToken, userId, os.Getenv("keycloak_realm"), body.Password, false)
	if err != nil {
		err2 := s.keycloak.DeleteUser(s.ctx, adminToken.AccessToken, os.Getenv("keycloak_realm"), userId)
		if err2 != nil {
			return errors.New(err.Error() + ";" + err2.Error())
		}

		return err
	}

	role, err := s.keycloak.GetRealmRole(s.ctx, adminToken.AccessToken, os.Getenv("keycloak_realm"), "user")
	if err != nil {
		err2 := s.keycloak.DeleteUser(s.ctx, adminToken.AccessToken, os.Getenv("keycloak_realm"), userId)
		if err2 != nil {
			return errors.New(err.Error() + ";" + err2.Error())
		}

		return err
	}

	err = s.keycloak.AddRealmRoleToUser(s.ctx, adminToken.AccessToken, os.Getenv("keycloak_realm"), userId, []gocloak.Role{*role})
	if err != nil {
		err2 := s.keycloak.DeleteUser(s.ctx, adminToken.AccessToken, os.Getenv("keycloak_realm"), userId)
		if err2 != nil {
			return errors.New(err.Error() + ";" + err2.Error())
		}

		return err
	}

	return nil
}

func (s *Service) loginAdmin() (*gocloak.JWT, error) {
	token, err := s.keycloak.LoginAdmin(s.ctx, os.Getenv("keycloak_admin_username"), os.Getenv("keycloak_admin_password"), os.Getenv("keycloak_realm"))
	if err != nil {
		return nil, err
	}
	return token, nil
}
