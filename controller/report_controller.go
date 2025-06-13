package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/user-service/dto"
	"github.com/swimresults/user-service/model"
	"github.com/swimresults/user-service/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func reportController() {
	router.GET("/report", getReports)

	router.GET("/report/subject-types", getSubjectTypes)

	router.POST("/report", addReport)
	router.POST("/report/submit", submitReport)

	router.DELETE("/report/:id", removeReport)
}

func getReports(c *gin.Context) {
	reports, err := service.GetReports()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, reports)
}

func getSubjectTypes(c *gin.Context) {
	types := model.GetReportSubjectTypes()
	c.IndentedJSON(http.StatusOK, types)
}

func submitReport(c *gin.Context) {
	var submission dto.UserReportSubmission
	if err := c.BindJSON(&submission); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var report model.UserReport

	if !submission.Anonymous {
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

		report.UserId = user.Identifier
	}

	report.Message = submission.Message
	report.SubjectId = submission.SubjectId
	report.SubjectType = submission.SubjectType

	newReport, err := service.AddReport(report)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newReport)
}

func addReport(c *gin.Context) {
	if failIfNotRoot(c) {
		return
	}

	var report model.UserReport
	if err := c.BindJSON(&report); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	report, err := service.AddReport(report)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, report)
}

func removeReport(c *gin.Context) {
	if failIfNotRoot(c) {
		return
	}

	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	err := service.RemoveReport(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
