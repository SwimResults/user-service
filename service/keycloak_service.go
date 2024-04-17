package service

import (
	"context"
	"errors"
	"github.com/Nerzal/gocloak/v13"
	"time"
)

func GetKeycloakUser(uuid string) (*gocloak.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if keycloak == nil {
		return nil, errors.New("no keycloak")
	}

	user, err := keycloak.Client.GetUserByID(ctx, keycloak.Token.AccessToken, keycloak.Realm, uuid)
	if err != nil {
		return nil, err
	}

	return user, nil
}
