package service

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
)

type MockKeycloak struct{}

func NewMockKeycloak() *MockKeycloak {
	return &MockKeycloak{}
}

func (m *MockKeycloak) Login(ctx context.Context, clientID, clientSecret, realm, username, password string) (*gocloak.JWT, error) {
	return &gocloak.JWT{
		AccessToken:      "mock-access-token",
		RefreshToken:     "mock-refresh-token",
		ExpiresIn:        3600,
		RefreshExpiresIn: 7200,
	}, nil
}

func (m *MockKeycloak) RefreshToken(ctx context.Context, refreshToken, clientID, clientSecret, realm string) (*gocloak.JWT, error) {
	return &gocloak.JWT{
		AccessToken:      "mock-access-token",
		RefreshToken:     "mock-refresh-token",
		ExpiresIn:        3600,
		RefreshExpiresIn: 7200,
	}, nil
}

func (m *MockKeycloak) CreateUser(ctx context.Context, token, realm string, user gocloak.User) (string, error) {
	return "mock-user-id", nil
}

func (m *MockKeycloak) SetPassword(ctx context.Context, token, userID, realm, password string, temporary bool) error {
	return nil
}

func (m *MockKeycloak) DeleteUser(ctx context.Context, token, realm, id string) error {
	return nil
}

func (m *MockKeycloak) GetRealmRole(ctx context.Context, token, realm, roleName string) (*gocloak.Role, error) {
	return &gocloak.Role{
		Name: gocloak.StringP(roleName),
	}, nil
}

func (m *MockKeycloak) CreateRealmRole(ctx context.Context, token, realm string, role gocloak.Role) (string, error) {
	return "mock-role-id", nil
}

func (m *MockKeycloak) AddRealmRoleToUser(ctx context.Context, token, realm, userID string, roles []gocloak.Role) error {
	return nil
}

func (m *MockKeycloak) LoginAdmin(ctx context.Context, username, password, realm string) (*gocloak.JWT, error) {
	return &gocloak.JWT{
		AccessToken: "mock-admin-token",
	}, nil
}
