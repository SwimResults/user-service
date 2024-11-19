package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/user-service/service"
	"net/http"
)

func notificationController() {

	router.POST("/notification/test/:device", sendTestNotification)

	router.OPTIONS("/notification/test", okay)
}

func sendTestNotification(c *gin.Context) {

	if failIfNotRoot(c) {
		return
	}

	device := c.Param("device")

	err := service.SendTestPushNotification(device)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
