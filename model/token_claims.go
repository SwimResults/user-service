package model

import "github.com/google/uuid"

type TokenClaims struct {
	Sub    uuid.UUID
	Scopes []string
}

func (tc *TokenClaims) IsRoot() bool {
	for _, scope := range tc.Scopes {
		if scope == "root" {
			return true
		}
	}

	return false
}
