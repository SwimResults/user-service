package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/user-service/dto"
	"github.com/swimresults/user-service/service"
	"io"
	"net/http"
)

func notificationController() {

	router.POST("/notification/test/:device", sendTestNotification)
	router.POST("/notification/:device", sendNotification)
	router.POST("/notification/meet/:meeting", sendNotificationForMeeting)

	router.POST("/notification/broadcast/:channel", sendBroadcast)
	router.POST("/notification/broadcast/meeting/:meeting", sendMeetingBroadcast)

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

	apnsId, body, status, err := service.SendPushNotification(device, request.Title, request.Subtitle, request.Message, request.InterruptionLevel)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(status, dto.NotificationResponseDto{
		ApnsId: apnsId,
		Body:   body,
	})
}

func sendNotificationForMeeting(c *gin.Context) {

	if failIfNotRoot(c) {
		return
	}

	meeting := c.Param("meeting")

	var request dto.MeetingNotificationRequestDto
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	users, notificationUsers, success, err := service.SendPushNotificationForMeeting(meeting, request)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, dto.NotificationsResponseDto{
		UserCount:             users,
		NotificationUserCount: notificationUsers,
		SuccessCount:          success,
	})
}

func sendBroadcast(c *gin.Context) {

	if failIfNotRoot(c) {
		return
	}

	channel := c.Param("channel")

	content, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	print("channel: " + channel)
	print("content: " + string(content))

	apnsRequestId, apnsUniqueId, body, status, err := service.SendPushBroadcast(channel, string(content))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(status, dto.BroadcastResponseDto{
		ApnsRequestId: apnsRequestId,
		ApnsUniqueId:  apnsUniqueId,
		Body:          body,
	})
}

func sendMeetingBroadcast(c *gin.Context) {

	if failIfNotRoot(c) {
		return
	}

	meeting := c.Param("meeting")

	content, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	print("meeting: " + meeting)
	print("content: " + string(content))

	apnsRequestId, apnsUniqueId, body, status, err := service.SendPushMeetingBroadcast(meeting, string(content))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(status, dto.BroadcastResponseDto{
		ApnsRequestId: apnsRequestId,
		ApnsUniqueId:  apnsUniqueId,
		Body:          body,
	})
}
