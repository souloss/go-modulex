package repository

import (
	"context"
	"errors"

	"github.com/souloss/go-clean-arch/entities"
	"github.com/souloss/go-clean-arch/internal/query"
)

type ArticleRepo struct {
	q *query.Query
}

// 确保实现了接口
var _ entities.ArticleRepository = (*ArticleRepo)(nil)

func NewArticleRepo(q *query.Query) *ArticleRepo {
	_ = q.Article.UnderlyingDB().AutoMigrate(&entities.Article{})
	return &ArticleRepo{q: q}
}

// GetById implements entities.ArticleRepository.
func (a *ArticleRepo) GetById(ctx context.Context, id string) (*entities.Article, error) {
	return a.q.Article.WithContext(ctx).Select(a.q.Author.ID.Eq(id)).First()
}

// UpdateById implements entities.ArticleRepository.
func (a *ArticleRepo) UpdateById(ctx context.Context, id string, r *entities.UpdateArticleRequest) (*entities.Article, error) {
	// 检查记录是否存在
	record, err := a.q.Article.WithContext(ctx).Where(a.q.Article.ID.Eq(id)).First()
	if err != nil {
		return nil, errors.New("record not found")
	}
	record.Title = r.Title
	record.Content = r.Content
	record.Author = r.Author
	// 更新
	_, err = a.q.Article.WithContext(ctx).Updates(record)
	return record, err
}

// GetAll implements ArticleRepository.
func (a *ArticleRepo) GetAll(ctx context.Context) ([]*entities.Article, error) {
	return a.q.Article.WithContext(ctx).Find()
}

// GetByID implements ArticleRepository.
func (a *ArticleRepo) GetByID(ctx context.Context, id string) (*entities.Article, error) {
	return a.q.Article.WithContext(ctx).Where(a.q.Article.ID.Eq(id)).First()
}

// Save implements ArticleRepository.
func (a *ArticleRepo) Save(ctx context.Context, ar *entities.Article) error {
	err := a.q.Article.WithContext(ctx).Create(ar)
	if err != nil {
		return err
	}
	return nil
}

// Delete implements ArticleRepository.
func (a *ArticleRepo) DeleteById(ctx context.Context, id string) error {
	_, err := a.q.Article.WithContext(ctx).Where(a.q.Article.ID.Eq(id)).Delete()
	return err
}
