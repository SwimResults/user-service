package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/user-service/model"
	"github.com/swimresults/user-service/service"
	"net/http"
)

func widgetController() {
	router.GET("/widget", getWidgets)
	router.GET("/dashboard", getUserDashboard)
	router.GET("/dashboard/default", getDefaultDashboard)

	router.POST("/dashboard", addUserDashboard)

	router.DELETE("/dashboard", removeUserDashboard)

	router.OPTIONS("/dashboard", okay)
}

func getWidgets(c *gin.Context) {
	widgets, err := service.GetWidgets()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, widgets)
}

func getUserDashboard(c *gin.Context) {
	claims, err1 := getClaimsFromAuthHeader(c)

	if err1 != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err1.Error()})
		return
	}

	dashboard, _, err := service.GetDashboardForUser(claims.Sub)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, dashboard)
}

func getDefaultDashboard(c *gin.Context) {
	dashboard, err := service.GetDefaultDashboard()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, dashboard)
}

func addUserDashboard(c *gin.Context) {
	claims, err1 := getClaimsFromAuthHeader(c)

	if err1 != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err1.Error()})
		return
	}

	var dashboard model.Dashboard
	if err := c.BindJSON(&dashboard); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	dashboard, created, err := service.SetUserDashboard(dashboard, claims.Sub)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if created {
		c.IndentedJSON(http.StatusCreated, dashboard)
	} else {
		c.IndentedJSON(http.StatusOK, dashboard)
	}

}

func removeUserDashboard(c *gin.Context) {
	claims, err1 := getClaimsFromAuthHeader(c)

	if err1 != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err1.Error()})
		return
	}

	err := service.RemoveUserDashboard(claims.Sub)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
