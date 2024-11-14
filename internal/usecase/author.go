package usecase

import (
	"context"

	"github.com/souloss/go-clean-arch/entities"
	"github.com/souloss/go-clean-arch/pkg/logger"
)

type AuthorUsecase struct {
	authorRespo entities.AuthorRepository
	logger      logger.FormatStrLogger
}

// 确保实现了接口
var _ entities.AuthorUsecase = (*AuthorUsecase)(nil)

// NewAuthorUsecase creates a new application service: AuthorService.
func NewAuthorUsecase(repo entities.AuthorRepository) *AuthorUsecase {
	return &AuthorUsecase{
		authorRespo: repo,
		logger:      logger.L().Named("author"),
	}
}

// DeleteById implements entities.AuthorUsecase.
func (a *AuthorUsecase) DeleteById(ctx context.Context, req *entities.AuthoryIDRequest) error {
	return a.authorRespo.DeleteById(ctx, req.ID)
}

// GetAll implements entities.AuthorUsecase.
func (a *AuthorUsecase) GetAll(ctx context.Context) ([]*entities.Author, error) {
	return a.authorRespo.GetAll(ctx)
}

// GetById implements entities.AuthorUsecase.
func (a *AuthorUsecase) GetById(ctx context.Context, req *entities.AuthoryIDRequest) (*entities.Author, error) {
	return a.authorRespo.GetById(ctx, req.ID)
}

// Save implements entities.AuthorUsecase.
func (a *AuthorUsecase) Save(ctx context.Context, req *entities.AuthoryNameRequest) error {
	err := a.authorRespo.Save(ctx, req.Name)
	if err != nil {
		a.logger.Error("Save author failed, err: %v", err)
		return err
	}
	return nil
}

// UpdateById implements entities.AuthorUsecase.
func (a *AuthorUsecase) UpdateById(ctx context.Context, req *entities.UpdateAuthorRequest) (*entities.Author, error) {
	return a.authorRespo.UpdateById(ctx, req.ID, req.Name)
}
