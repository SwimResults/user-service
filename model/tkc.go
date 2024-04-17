package model

import "github.com/Nerzal/gocloak/v13"

type Tkc struct {
	Realm  string
	Token  *gocloak.JWT
	Client *gocloak.GoCloak
}
