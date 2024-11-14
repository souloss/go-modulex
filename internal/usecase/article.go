package usecase

import (
	"context"

	"github.com/souloss/go-clean-arch/entities"
	"github.com/souloss/go-clean-arch/pkg/logger"
)

type ArticleUsecase struct {
	articleRespo entities.ArticleRepository
	logger       logger.FormatStrLogger
}

// 确保实现了接口
var _ entities.ArticleUsecase = (*ArticleUsecase)(nil)

// NewArticleUsecase creates a new application service: ArticleService.
func NewArticleUsecase(repo entities.ArticleRepository) *ArticleUsecase {
	return &ArticleUsecase{
		articleRespo: repo,
		logger:       logger.L().Named("article"),
	}
}

// DeleteById implements entities.ArticleUsecase.
func (s *ArticleUsecase) DeleteById(ctx context.Context, req *entities.ArticleByIDRequest) error {
	return s.articleRespo.DeleteById(ctx, req.ID)
}

// GetAll implements entities.ArticleUsecase.
func (s *ArticleUsecase) GetAll(ctx context.Context) ([]*entities.Article, error) {
	return s.articleRespo.GetAll(ctx)
}

// GetById implements entities.ArticleUsecase.
func (s *ArticleUsecase) GetById(ctx context.Context, req *entities.ArticleByIDRequest) (*entities.Article, error) {
	return s.articleRespo.GetById(ctx, req.ID)
}

// Save implements entities.ArticleUsecase.
func (s *ArticleUsecase) Save(ctx context.Context, r *entities.SaveArticleRequest) error {
	return s.articleRespo.Save(ctx, &entities.Article{
		Title:   r.Title,
		Content: r.Content,
		Author:  r.Author,
	})
}

// UpdateById implements entities.ArticleUsecase.
func (s *ArticleUsecase) UpdateById(ctx context.Context, r *entities.UpdateArticleRequest) (*entities.Article, error) {
	return s.articleRespo.UpdateById(ctx, r.ID, r)
}
