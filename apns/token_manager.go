package apns

import (
	"encoding/json"
	"fmt"
	"github.com/swimresults/user-service/model"
	"os"
	"sync"
)

var token Token

func GetToken() string {
	token.GenerateIfExpired()
	return token.Bearer
}

func Init() {

	authKey, err := AuthKeyFromFile("config/apns/AuthKey.p8")
	if err != nil {
		println("failed reading apns token file")
	}

	dat, err1 := os.ReadFile("config/apns/token_config.json")
	if err1 != nil {
		println(err1.Error())
		return
	}

	var tokenConfig model.ApnsTokenConfig

	err = json.Unmarshal(dat, &tokenConfig)
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Printf("set token config to: key: '%s'; team: '%s'\n", tokenConfig.KeyId, tokenConfig.TeamId)

	token = Token{
		Mutex:   sync.Mutex{},
		AuthKey: authKey,
		KeyID:   tokenConfig.KeyId,
		TeamID:  tokenConfig.TeamId,
	}
}
