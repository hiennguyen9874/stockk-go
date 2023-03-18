package usecase

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/config"
	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/studytemplates"
	"github.com/hiennguyen9874/stockk-go/internal/usecase"
	"github.com/hiennguyen9874/stockk-go/pkg/logger"
)

type studyTemplateUseCase struct {
	usecase.UseCase[models.StudyTemplate]
	studyTemplatePgRepo studytemplates.StudyTemplatePgRepository
}

func StudyTemplateUseCaseI(
	studyTemplatePgRepo studytemplates.StudyTemplatePgRepository,
	cfg *config.Config,
	logger logger.Logger,
) studytemplates.StudyTemplateUseCaseI {
	return &studyTemplateUseCase{
		UseCase:             usecase.CreateUseCase[models.StudyTemplate](studyTemplatePgRepo, cfg, logger),
		studyTemplatePgRepo: studyTemplatePgRepo,
	}
}

func (u *studyTemplateUseCase) GetByOwnerName(ctx context.Context, ownerSource string, ownerId string, name string) (*models.StudyTemplate, error) {
	return u.studyTemplatePgRepo.GetByOwnerName(ctx, ownerSource, ownerId, name)
}

func (u *studyTemplateUseCase) GetOrCreateWithOwnerName(ctx context.Context, ownerSource string, ownerId string, name string, content string) (*models.StudyTemplate, bool, error) {
	return u.studyTemplatePgRepo.GetOrCreateWithOwnerName(ctx, ownerSource, ownerId, name, content)
}

func (u *studyTemplateUseCase) GetAllByOwner(ctx context.Context, ownerSource string, ownerId string) ([]*models.StudyTemplate, error) {
	return u.studyTemplatePgRepo.GetAllByOwner(ctx, ownerSource, ownerId)
}

func (u *studyTemplateUseCase) CreateOrUpdateWithOwnerName(ctx context.Context, ownerSource string, ownerId string, name string, content string) (*models.StudyTemplate, bool, bool, error) {
	studyTemplate, created, err := u.studyTemplatePgRepo.GetOrCreateWithOwnerName(ctx, ownerSource, ownerId, name, content)

	if err != nil {
		return nil, false, false, err
	}

	if created {
		return studyTemplate, true, false, nil
	}

	values := make(map[string]interface{})
	values["content"] = content

	updatedStudyTemplate, err := u.studyTemplatePgRepo.Update(ctx, studyTemplate, values)
	if err != nil {
		return nil, false, false, err
	}

	return updatedStudyTemplate, false, true, nil
}

func (u *studyTemplateUseCase) DeleteByOwnerName(ctx context.Context, ownerSource string, ownerId string, name string) (*models.StudyTemplate, error) {
	return u.studyTemplatePgRepo.DeleteByOwnerName(ctx, ownerSource, ownerId, name)
}
