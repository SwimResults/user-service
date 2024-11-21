package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/user-service/dto"
	"github.com/swimresults/user-service/service"
	"net/http"
)

func notificationController() {

	router.POST("/notification/test/:device", sendTestNotification)
	router.POST("/notification/:device", sendNotification)

	router.OPTIONS("/notification/test/:device", okay)
	router.OPTIONS("/notification/:device", okay)
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

func sendNotification(c *gin.Context) {

	if failIfNotRoot(c) {
		return
	}

	device := c.Param("device")

	var request dto.NotificationRequestDto
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	apnsId, body, status, err := service.SendPushNotification(device, request.Title, request.Subtitle, request.Message)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(status, dto.NotificationResponseDto{
		ApnsId: apnsId,
		Body:   body,
	})
}
