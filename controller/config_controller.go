package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/user-service/service"
	"net/http"
)

func configController() {
	router.GET("/config", getConfigs)

	router.POST("/config/enable/:meeting/:enable", setEnable)
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
