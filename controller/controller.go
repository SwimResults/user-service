package controller

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/swimresults/user-service/model"
	"github.com/swimresults/user-service/service"
	"net/http"
	"os"
	"strings"
)

var router = gin.Default()
var serviceKey string

func Run() {

	port := os.Getenv("SR_USER_PORT")

	if port == "" {
		fmt.Println("no application port given! Please set SR_USER_PORT.")
		return
	}

	serviceKey = os.Getenv("SR_SERVICE_KEY")

	if serviceKey == "" {
		fmt.Println("no security for inter-service communication given! Please set SR_SERVICE_KEY.")
		return
	}

	userController()
	widgetController()
	dashboardController()
	notificationUserController()
	notificationController()

	router.GET("/actuator", actuator)

	err := router.Run(":" + port)
	if err != nil {
		fmt.Println("Unable to start application on port " + port)
		return
	}
}

func actuator(c *gin.Context) {

	state := "OPERATIONAL"

	if !service.PingDatabase() {
		state = "DATABASE_DISCONNECTED"
	}
	c.String(http.StatusOK, state)
}

func checkServiceKey(c *gin.Context) error {
	if c.Request.Header["X-SWIMRESULTS-SERVICE"][0] == serviceKey {
		return nil
	}

	return errors.New("no service authorization key in header")
}

func checkAuthHeaderToken(c *gin.Context) error {
	claims, err1 := getClaimsFromAuthHeader(c)

	if err1 != nil {
		return err1
	}

	if !claims.IsRoot() {
		return errors.New("insufficient permissions")
	}

	return nil
}

func getClaimsFromAuthHeader(c *gin.Context) (*model.TokenClaims, error) {
	if len(c.Request.Header["Authorization"]) == 0 {
		err1 := errors.New("no authorization in header")
		c.IndentedJSON(http.StatusUnauthorized, err1.Error())
		return nil, err1
	}

	tokenString := strings.Split(c.Request.Header["Authorization"][0], " ")[1]

	token, err1 := jwt.Parse(tokenString, nil)
	if token == nil {
		return nil, err1
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	sub := fmt.Sprintf("%s", claims["sub"])

	id, err2 := uuid.Parse(sub)
	if err2 != nil {
		return nil, err2
	}

	var tokenClaims model.TokenClaims

	tokenClaims.Sub = id
	tokenClaims.Scopes = strings.Split(fmt.Sprintf("%s", claims["scope"]), " ")

	return &tokenClaims, nil
}

func checkIfRoot(c *gin.Context) error {
	if checkServiceKey(c) == nil {
		return nil
	}

	tokenError := checkAuthHeaderToken(c)

	if tokenError == nil {
		return nil
	} else {
		return tokenError
	}
}

// failIfNotRoot returns true if the requester is not root or a service
func failIfNotRoot(c *gin.Context) bool {
	err := checkIfRoot(c)

	if err == nil {
		return false
	} else {
		c.IndentedJSON(http.StatusUnauthorized, err.Error())
		return true
	}
}
