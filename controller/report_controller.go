package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swimresults/service-core/security"
	"github.com/swimresults/user-service/dto"
	"github.com/swimresults/user-service/model"
	"github.com/swimresults/user-service/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func reportController() {
	security.Route(router, "GET", "/report", security.PermissionAdmin, getReports)

	security.Route(router, "GET", "/report/subject-types", security.PermissionPublic, getSubjectTypes)
	security.Route(router, "GET", "/report/meet/:meeting", security.PermissionAdmin, getReportsByMeeting)

	security.Route(router, "POST", "/report", security.PermissionAdmin, addReport)
	security.Route(router, "POST", "/report/submit", security.PermissionPublic, submitReport)

	security.Route(router, "POST", "/report/:id/acknowledge", security.PermissionAdmin, acknowledgeReport)
	security.Route(router, "POST", "/report/:id/complete", security.PermissionAdmin, completeReport)

	security.Route(router, "DELETE", "/report/:id", security.PermissionAdmin, removeReport)
}

func getReports(c *gin.Context) {
	reports, err := service.GetReports()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, reports)
}

func getReportsByMeeting(c *gin.Context) {
	meeting := c.Param("meeting")
	if meeting == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "no meeting given"})
		return
	}

	reports, err := service.GetReportsByMeeting(meeting)
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
	report.Meeting = submission.Meeting

	newReport, err := service.AddReport(report)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newReport)
}

func addReport(c *gin.Context) {
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

func acknowledgeReport(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	report, err := service.ToggleReportAcknowledged(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, report)
}

func completeReport(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	report, err := service.ToggleReportComplete(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, report)
}

func removeReport(c *gin.Context) {
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
