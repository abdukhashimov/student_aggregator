package services

import (
	"context"
	"errors"
	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/internal/core/ports"
	"github.com/abdukhashimov/student_aggregator/pkg/parser"
)

var _ ports.AggregatorService = (*AggregatorService)(nil)

type AggregatorService struct {
	studentsRepo ports.StudentsStore
	schemasRepo  ports.SchemaStore
	storage      ports.StorageService
}

func NewAggregatorService(studentsRepo ports.StudentsStore, schemasRepo ports.SchemaStore, storage ports.StorageService) *AggregatorService {
	return &AggregatorService{
		studentsRepo: studentsRepo,
		schemasRepo:  schemasRepo,
		storage:      storage,
	}
}

func (aggS *AggregatorService) ParseFile(ctx context.Context, fileName string, schemaID string) error {
	schema, err := aggS.schemasRepo.GetById(ctx, schemaID)
	if err != nil {
		return err
	}

	schemaName := ""
	if schema.Name == domain.RSS {
		schemaName = domain.RSS
	}
	if schema.Name == domain.WAC {
		schemaName = domain.WAC
	}

	if schemaName == "" {
		return errors.New("unsupported schema")
	}

	r, _, err := aggS.storage.GetFile(ctx, fileName)
	if err != nil {
		return err
	}

	students := new([]domain.StudentRecord)
	err = parser.ParseXLSXFile(students, r, schema.ConvertToParserSchema())
	if err != nil {
		return err
	}

	// TODO: use transactions to avoid data inconsistency
	switch schemaName {
	case domain.RSS:
		{
			for _, student := range *students {
				_, sErr := aggS.studentsRepo.SaveRSS(ctx, fileName, student.Email, student.StudentRSS)
				if sErr != nil {
					return sErr
				}
			}
		}
	case domain.WAC:
		{
			for _, student := range *students {
				_, sErr := aggS.studentsRepo.SaveWAC(ctx, fileName, student.Email, student.StudentWAC)
				if sErr != nil {
					return sErr
				}
			}
		}
	}

	return nil
}
