package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/swimresults/user-service/dto"
	"github.com/swimresults/user-service/model"
	"github.com/swimresults/user-service/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

func userController() {
	router.GET("/users", getUsers)
	router.GET("/user", getUser)
	router.GET("/user/:id", getUserById)

	router.POST("/user", addUser)
	router.POST("/user/athlete", changeFollowerForUser)

	router.DELETE("/user/:id", removeUser)

	router.PUT("/user", updateUser)
}

func getUsers(c *gin.Context) {
	users, err := service.GetUsers()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {

	id, err1 := getUUIDFromRequest(c)

	if err1 != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err1.Error()})
		return
	}

	user, err := service.GetUserByKeycloakId(*id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func getUserById(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	user, err := service.GetUserById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func removeUser(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	err := service.RemoveUserById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusNoContent, "")
}

func addUser(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := service.AddUser(user)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func changeFollowerForUser(c *gin.Context) {
	id, err1 := getUUIDFromRequest(c)

	if err1 != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err1.Error()})
		return
	}

	var request dto.FollowAthleteRequestDto
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	user, err2 := service.ModifyFollowForUser(*id, request.AthleteId, request.Follow)
	if err2 != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func updateUser(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := service.UpdateUser(user)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func getUUIDFromRequest(c *gin.Context) (*uuid.UUID, error) {
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

	return &id, nil
}
