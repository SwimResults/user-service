package controller

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/swimresults/service-core/security"
	"github.com/swimresults/user-service/model"
	"github.com/swimresults/user-service/service"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var router = gin.Default()
var serviceKey string

func Run() {

	port := os.Getenv("SR_USER_PORT")

	if port == "" {
		fmt.Println("no application port given! Please set SR_USER_PORT.")
		return
	}

	security.InitAuthMiddleware(&security.AuthMiddlewareConfig{
		ServiceKey:    os.Getenv("SR_SERVICE_KEY"),
		ExcludedPaths: []string{"/actuator"},
	})

	p := ginprometheus.NewWithConfig(ginprometheus.Config{
		Subsystem: "gin",
	})
	p.Use(router)

	router.Use(security.AuthMiddleware())

	userController()
	widgetController()
	dashboardController()
	notificationUserController()
	notificationController()
	configController()
	reportController()

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
	println("checking service key...")
	received := c.Request.Header["X-Swimresults-Service"]
	fmt.Printf("received: '%s', expected: '%s'\n", received, serviceKey)
	if len(received) <= 0 {
		return errors.New("no service authorization key in header")
	}
	if received[0] == serviceKey {
		return nil
	}

	return errors.New("invalid service authorization key in header")
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
	claims, err1 := security.ValidateAuthorizationHeader(c.GetHeader("Authorization"))
	if err1 != nil {
		c.IndentedJSON(http.StatusUnauthorized, err1.Error())
		return nil, err1
	}

	sub := fmt.Sprintf("%s", claims.Subject)

	id, err2 := uuid.Parse(sub)
	if err2 != nil {
		return nil, err2
	}

	var tokenClaims model.TokenClaims

	tokenClaims.Sub = id
	tokenClaims.Scopes = strings.Fields(claims.Scope)

	return &tokenClaims, nil
}

func checkIfRoot(c *gin.Context) error {
	keyError := checkServiceKey(c)
	if keyError == nil {
		return nil
	}

	tokenError := checkAuthHeaderToken(c)

	if tokenError == nil {
		return nil
	} else {
		fmt.Printf("both auth checks for root failed: \n%s\n%s\n", keyError, tokenError)
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
