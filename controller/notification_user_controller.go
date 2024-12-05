package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/user-service/dto"
	"github.com/swimresults/user-service/model"
	"github.com/swimresults/user-service/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func notificationUserController() {
	router.GET("/notification_users", getNotificationUsers)
	router.GET("/notification_user", getNotificationUser)
	router.GET("/notification_user/token/:token", getNotificationUserByToken)
	router.GET("/notification_user/:id", getNotificationUserById)

	router.POST("/notification_user", addNotificationUser)
	router.POST("/notification_user/register", registerNotificationUserWithoutToken)
	router.POST("/notification_user/register/user", registerNotificationUser)

	router.DELETE("/notification_user/:id", removeNotificationUser)

	router.PUT("/notification_user", updateNotificationUser)

	router.OPTIONS("/notification_user", okay)
	router.OPTIONS("/notification_user/register", okay)
	router.OPTIONS("/notification_user/register/user", okay)
}

func getNotificationUsers(c *gin.Context) {

	if failIfNotRoot(c) {
		return
	}

	users, err := service.GetNotificationUsers()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

func getNotificationUser(c *gin.Context) {

	claims, err1 := getClaimsFromAuthHeader(c)

	if err1 != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err1.Error()})
		return
	}

	user, err := service.GetUserByKeycloakId(claims.Sub)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	notificationUser, err2 := service.GetNotificationUserByUserId(user.Identifier)
	if err2 != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err2.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, notificationUser)
}

func getNotificationUserByToken(c *gin.Context) {

	token := c.Param("token")

	if token == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no token given"})
		return
	}

	notificationUser, err2 := service.GetNotificationUserByToken(token)
	if err2 != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err2.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, notificationUser)
}

func getNotificationUserById(c *gin.Context) {

	if failIfNotRoot(c) {
		return
	}

	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	user, err := service.GetNotificationUserById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func removeNotificationUser(c *gin.Context) {

	if failIfNotRoot(c) {
		return
	}

	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	err := service.RemoveNotificationUserById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusNoContent, "")
}

func addNotificationUser(c *gin.Context) {

	if failIfNotRoot(c) {
		return
	}

	var user model.NotificationUser
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := service.AddNotificationUser(user)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func updateNotificationUser(c *gin.Context) {

	if failIfNotRoot(c) {
		return
	}

	var user model.NotificationUser
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := service.UpdateNotificationUser(user)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func registerNotificationUserWithoutToken(c *gin.Context) {
	var request dto.RegisterNotificationUserRequestDto
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := service.RegisterNotificationUser(request.Token, request.Device, request.Settings, nil)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func registerNotificationUser(c *gin.Context) {
	var request dto.RegisterNotificationUserRequestDto
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	claims, err1 := getClaimsFromAuthHeader(c)

	var user model.User
	if err1 == nil {
		user, err1 = service.GetUserByKeycloakId(claims.Sub)
		if err1 != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err1.Error()})
			return
		}
	}

	r, err := service.RegisterNotificationUser(request.Token, request.Device, request.Settings, &user)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}
