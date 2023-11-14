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

	router.POST("/widget", addWidget)

	router.DELETE("/widget/:id", removeWidget)
}

func getWidgets(c *gin.Context) {
	widgets, err := service.GetWidgets()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, widgets)
}

func addWidget(c *gin.Context) {
	if failIfNotRoot(c) {
		return
	}

	var widget model.Widget
	if err := c.BindJSON(&widget); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	widget, err := service.AddWidget(widget)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, widget)
}

func removeWidget(c *gin.Context) {
	if failIfNotRoot(c) {
		return
	}

	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	err := service.RemoveWidget(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
