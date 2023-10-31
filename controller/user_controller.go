package controller

import (
	"errors"
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
	router.POST("/user/me", changeMe)
	router.POST("/user/language", updateUserLanguage)

	router.DELETE("/user/:id", removeUser)

	router.PUT("/user", updateUser)

	router.OPTIONS("/user", okay)
	router.OPTIONS("/user/athlete", okay)
}

func okay(c *gin.Context) {
	c.Status(200)
}

func getUsers(c *gin.Context) {

	if failIfNotRoot(c) {
		return
	}

	users, err := service.GetUsers()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {

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

	c.IndentedJSON(http.StatusOK, user)
}

func getUserById(c *gin.Context) {

	if failIfNotRoot(c) {
		return
	}

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

	if failIfNotRoot(c) {
		return
	}

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

	if failIfNotRoot(c) {
		return
	}

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
	claims, err1 := getClaimsFromAuthHeader(c)

	if err1 != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err1.Error()})
		return
	}

	var request dto.FollowAthleteRequestDto
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	user, err2 := service.ModifyFollowForUser(claims.Sub, request.AthleteId, request.Follow)
	if err2 != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func changeMe(c *gin.Context) {
	claims, err1 := getClaimsFromAuthHeader(c)

	if err1 != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err1.Error()})
		return
	}

	var request dto.SetMeRequestDto
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	user, err2 := service.ModifyMe(claims.Sub, request.AthleteId, request.Set)
	if err2 != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func updateUserLanguage(c *gin.Context) {
	claims, err1 := getClaimsFromAuthHeader(c)

	if err1 != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err1.Error()})
		return
	}

	var request dto.SetLanguageRequestDto
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	user, err2 := service.ModifyUserLanguage(claims.Sub, request.Language)
	if err2 != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func updateUser(c *gin.Context) {

	if failIfNotRoot(c) {
		return
	}

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

func getClaimsFromAuthHeader(c *gin.Context) (*model.TokenClaims, error) {
	if len(c.Request.Header["Authorization"]) == 0 {
		err1 := errors.New("no authorization in header")
		c.IndentedJSON(http.StatusUnauthorized, err1.Error())
		return nil, err1
	}

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

	var tokenClaims model.TokenClaims

	tokenClaims.Sub = id
	tokenClaims.Scopes = strings.Split(fmt.Sprintf("%s", claims["scope"]), " ")

	return &tokenClaims, nil
}

func failIfNotRoot(c *gin.Context) bool {

	claims, err1 := getClaimsFromAuthHeader(c)

	if err1 != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err1.Error()})
		return true
	}

	if !claims.IsRoot() {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "insufficient permissions"})
		return true
	}

	return false
}
