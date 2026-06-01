package model

import "github.com/google/uuid"

type TokenClaims struct {
	Sub    uuid.UUID
	Scopes []string
}

func (tc *TokenClaims) IsAdmin() bool {
	for _, scope := range tc.Scopes {
		if scope == "admin" {
			return true
		}
	}

	return false
}
