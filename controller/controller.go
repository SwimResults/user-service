package controller

import (
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
