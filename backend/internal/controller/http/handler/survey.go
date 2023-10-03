package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mado/internal/core/survey"
)

type SurveyService interface {
	// Create(ctx context.Context, dto user.CreateDTO) (user.User, error)
	Create(*survey.SurveyRequirements) (*survey.SurveyRequirements, error)
	GetSurviesByUserID(user_id int, ctx *gin.Context) (response []*survey.SurveyResponse, err error)
	// GetAllRows()()
}

type surveyDeps struct {
	router *gin.RouterGroup

	surveyService SurveyService
}

type surveyHandler struct {
	surveyService SurveyService
}

func newSurveyHandler(deps surveyDeps) {
	handler := surveyHandler{
		surveyService: deps.surveyService,
	}

	usersGroup := deps.router.Group("/survey")
	{
		usersGroup.POST("/create", handler.CreateSurvey)
		usersGroup.GET("/get/:user_id", handler.GetSurveis)
	}

}

func (h surveyHandler) CreateSurvey(c *gin.Context) {
	var request *survey.SurveyRequirements
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	resp, err := h.surveyService.Create(request)
	if err != nil {
		switch err {
		case survey.ErrSurvey:
			c.JSON(http.StatusBadRequest, gin.H{"error": "SurveyRequirements is nil"})
		case survey.ErrSurveyName:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Name field is empty"})
		case survey.ErrSurveyQuestion:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Questions field is empty"})
		default:
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"created_survey": resp})

}

func (h surveyHandler) GetSurveis(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))

	respnose, err := h.surveyService.GetSurviesByUserID(userID, c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, respnose)
}
