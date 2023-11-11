package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/user-service/model"
	"github.com/swimresults/user-service/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func widgetController() {
	router.GET("/widget", getWidgets)
	router.GET("/dashboard", getUserDashboard)
	router.GET("/dashboard/default", getDefaultDashboard)

	router.POST("/dashboard", addUserDashboard)

	router.DELETE("/dashboard/:id", removeUserDashboard)

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

	state := c.Query("meeting_state")

	dashboard, _, err := service.GetDashboardForUser(state, claims.Sub)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, dashboard)
}

func getDefaultDashboard(c *gin.Context) {
	state := c.Query("meeting_state")
	dashboard, err := service.GetDefaultDashboard(state)
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

	dashboard, err := service.AddUserDashboard(dashboard, claims.Sub)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, dashboard)
}

func removeUserDashboard(c *gin.Context) {
	claims, err1 := getClaimsFromAuthHeader(c)

	if err1 != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err1.Error()})
		return
	}

	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	err := service.RemoveUserDashboard(id, claims.Sub)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
