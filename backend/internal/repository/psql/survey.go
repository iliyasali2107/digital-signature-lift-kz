package psql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"mado/internal/core/survey"
	"mado/pkg/database/postgres"
	"mado/pkg/logger"
)

// Survey is a Survey repository.
type SurveyrRepository struct {
	db     *postgres.Postgres
	logger *zap.Logger
}

// NewSurveyRepository creates a new UserRepository.
func NewSurveyrRepository(db *postgres.Postgres, logger *zap.Logger) SurveyrRepository {
	return SurveyrRepository{
		db:     db,
		logger: logger,
	}
}

func (s SurveyrRepository) Create(req *survey.SurveyRequirements, ctx context.Context) (*survey.SurveyRequirements, error) {
	tx, err := s.startTransaction(ctx)
	if err != nil {
		s.logger.Error("error in starting transaction: ", zap.Error(err))
		return nil, err
	}

	defer s.rollbackIfError(tx, ctx, &err)

	questionIDs, err := s.insertQuestions(tx, ctx, req.Questions)
	if err != nil {
		s.logger.Error("error in inserting questions for survey: ", zap.Error(err))
		return nil, err
	}

	if err := s.insertSurvey(tx, ctx, req, questionIDs); err != nil {
		s.logger.Error("error in inserting survey: ", zap.Error(err))
		return nil, err
	}

	if err := s.commitTransaction(tx, ctx); err != nil {
		s.logger.Error("error in commiting transaction: ", zap.Error(err))
		return nil, err
	}

	return req, nil
}

func (s SurveyrRepository) GetSurviesByUserID(user_id int, ctx *gin.Context) (response []*survey.SurveyResponse, err error) {
	query := "SELECT * FROM survey WHERE user_id = $1"
	rows, err := s.db.Pool.Query(ctx, query, user_id)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()
	// var surveyResponse *survey.SurveyResponse
	var surveis []*survey.SurveyResponse
	surveyResponse := new(survey.SurveyResponse)
	for rows.Next() {
		// Scan the row into variables
		if err := rows.Scan(&surveyResponse.ID, &surveyResponse.Name, &surveyResponse.Status, &surveyResponse.Rka, &surveyResponse.Rc_name, &surveyResponse.Adress, &surveyResponse.Question_id, &surveyResponse.CreatedAt, &surveyResponse.User_id); err != nil {
			s.logger.Error("GetSurviesByUserID scanning err: ", zap.Error(err))
			log.Fatalf("Error scanning row: %v", err)
		}
		surveis = append(surveis, surveyResponse)
	}
	fmt.Println("surveyResponse:", surveis)
	return surveis, nil
}

func (s SurveyrRepository) startTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.db.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}
	return tx, nil
}

func (s SurveyrRepository) rollbackIfError(tx pgx.Tx, ctx context.Context, err *error) {
	if *err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			fmt.Printf("Error rolling back transaction: %v\n", rollbackErr)
		}
	}
}

func (s SurveyrRepository) insertQuestions(tx pgx.Tx, ctx context.Context, questions []survey.Question) ([]int, error) {
	questionIDs := []int{}
	for _, q := range questions {
		var questionID int
		err := tx.QueryRow(ctx, "INSERT INTO public.question (description) VALUES ($1) RETURNING id", q.Description).Scan(&questionID)
		if err != nil {
			return nil, fmt.Errorf("could not insert question: %w", err)
		}
		questionIDs = append(questionIDs, questionID)
	}
	//todo delete this
	// Convert questionIDs to a JSON array string
	questionIDsJSON, err := json.Marshal(questionIDs)
	if err != nil {
		return nil, fmt.Errorf("could not marshal question IDs: %w", err)
	}

	s.logger.Info("Questions of survey", zap.String("question_ids", string(questionIDsJSON)))
	return questionIDs, nil
}

func (s SurveyrRepository) insertSurvey(tx pgx.Tx, ctx context.Context, req *survey.SurveyRequirements, questionIDs []int) error {
	//todo instead of mocks use request's value
	mockRka := "Mock Rka Value"
	mockRcName := "Mock RcName Value"
	mockAdress := "Mock Address Value"

	insertBuilder := s.db.Builder.Insert("public.survey").
		Columns("name", "rka", "rc_name", "adress", "question_id", "user_id").
		Values(req.Name, mockRka, mockRcName, mockAdress, questionIDs, req.UserID).
		Suffix("RETURNING id")

	sqlQuery, args, err := insertBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("can not build insert survey query: %w", err)
	}

	//todo maybe del it
	// Convert args to a JSON string
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("can not marshal args: %w", err)
	}
	s.logger.Info("InsertSurvey query", zap.String("sql", sqlQuery), zap.String("args", string(argsJSON)))

	logger.FromContext(ctx).Debug("check following query", zap.String("sql", sqlQuery), zap.Any("args", args))

	var id string
	if err := s.db.Pool.QueryRow(ctx, sqlQuery, args...).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return fmt.Errorf("can not insert user: %w", survey.ErrAlreadyExist)
			}
		}
		return fmt.Errorf("can not insert survey: %w", err)
	}
	req.ID = id
	return nil
}

func (s SurveyrRepository) commitTransaction(tx pgx.Tx, ctx context.Context) error {
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}
	return nil

}
