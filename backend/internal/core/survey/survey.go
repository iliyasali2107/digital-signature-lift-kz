package survey

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Repository is a user repository.
type Repository interface {
	// Create(*survey.SurveyRequirements) (*survey.SurveyRequirements, error)
	Create(*SurveyRequirements, context.Context) (*SurveyRequirements, error)
	GetSurviesByUserID(user_id int, ctx *gin.Context) (response []*SurveyResponse, err error)
}

// Service is a user service interface.
type Service struct {
	surveyRepository Repository
	logger           *zap.Logger
}

// NewService creates a new user service.
func NewService(surveyRepository Repository, logger *zap.Logger) Service {
	return Service{
		surveyRepository: surveyRepository,
		logger:           logger,
	}
}

func (s Service) Create(requirements *SurveyRequirements) (*SurveyRequirements, error) {

	if err := s.ValidateSurveyRequirements(requirements); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.surveyRepository.Create(requirements, ctx)
}
func (s Service) GetSurviesByUserID(user_id int, ctx *gin.Context) (response []*SurveyResponse, err error) {
	return s.surveyRepository.GetSurviesByUserID(user_id, ctx)
}
func (s Service) ValidateSurveyRequirements(requirements *SurveyRequirements) error {
	if requirements == nil {
		return ErrSurvey
	}

	if requirements.Name == "" {
		return ErrSurveyName
	}

	if len(requirements.Questions) == 0 {
		return ErrSurveyQuestion
	}

	// Add more checks for other fields as needed

	return nil
}
