package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/user-service/model"
	"github.com/swimresults/user-service/service"
	"net/http"
)

func configController() {
	router.GET("/config", getConfigs)

	router.POST("/config/enable/:meeting/:enable", setEnable)
	router.POST("/config", addConfig)
}

func getConfigs(c *gin.Context) {

	if failIfNotRoot(c) {
		return
	}

	configs, err := service.GetConfigs()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, configs)
}

func setEnable(c *gin.Context) {

	if failIfNotRoot(c) {
		return
	}

	meeting := c.Param("meeting")
	if meeting == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no meeting given"})
		return
	}

	enable := c.Param("enable")

	config, err2 := service.DisableNotification(meeting, enable == "true")
	if err2 != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err2.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, config)
}

func addConfig(c *gin.Context) {

	if failIfNotRoot(c) {
		return
	}

	var config model.Config
	if err := c.BindJSON(&config); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := service.AddConfig(config)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}
