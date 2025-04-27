package usecase

import (
	"xyz_multifinance/src/internal/source/model"
	"xyz_multifinance/src/internal/source/repository"

	"golang.org/x/net/context"
)

type Usecase struct {
	repo repository.Repository
}

func NewUsecase(repo repository.Repository) Usecase {
	return Usecase{repo: repo}
}

type Source struct {
	ID       string
	Secret   string
	Category string
	Name     string
	Email    string
}

func (u Usecase) Register(ctx context.Context, request Source) (*model.Source, error) {
	source := &model.Source{
		ID:         request.ID,
		SecretHash: "",
		Category:   request.Category,
		Name:       request.Name,
		Email:      request.Email,
	}

	if err := source.GenerateHashFromSecret(request.Secret); err != nil {
		return nil, err
	}

	if err := u.repo.Create(ctx, source); err != nil {
		return nil, err
	}

	return source, nil

}

func (u Usecase) FindByID(ctx context.Context, sourceID string) (*model.Source, error) {
	return u.repo.FindByID(ctx, sourceID)
}
