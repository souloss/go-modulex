package repository

import (
	"context"
	"errors"

	"github.com/souloss/go-clean-arch/entities"
	"github.com/souloss/go-clean-arch/internal/query"
)

type AuthorRepo struct {
	q *query.Query
}

// 确保实现了接口
var _ entities.AuthorRepository = (*AuthorRepo)(nil)

func NewAuthorRepo(q *query.Query) *AuthorRepo {
	_ = q.Author.UnderlyingDB().AutoMigrate(&entities.Author{})
	return &AuthorRepo{q: q}
}

// GetById implements entities.ArticleRepository.
func (a *AuthorRepo) GetById(ctx context.Context, id string) (*entities.Author, error) {
	return a.q.Author.WithContext(ctx).Select(a.q.Author.ID.Eq(id)).First()
}

// UpdateById implements entities.ArticleRepository.
func (a *AuthorRepo) UpdateById(ctx context.Context, id string, name string) (*entities.Author, error) {
	// 检查记录是否存在
	record, err := a.q.Author.WithContext(ctx).Where(a.q.Author.ID.Eq(id)).First()
	if err != nil {
		return nil, errors.New("record not found")
	}
	record.Name = name
	// 更新
	_, err = a.q.Author.WithContext(ctx).Updates(record)
	return record, err
}

// GetAll implements ArticleRepository.
func (a *AuthorRepo) GetAll(ctx context.Context) ([]*entities.Author, error) {
	return a.q.Author.WithContext(ctx).Find()
}

// GetByID implements ArticleRepository.
func (a *AuthorRepo) GetByID(ctx context.Context, id string) (*entities.Author, error) {
	return a.q.Author.WithContext(ctx).Where(a.q.Author.ID.Eq(id)).First()
}

// Save implements ArticleRepository.
func (a *AuthorRepo) Save(ctx context.Context, name string) error {
	err := a.q.Author.WithContext(ctx).Create(&entities.Author{
		Name: name,
	})
	if err != nil {
		return err
	}
	return nil
}

// Delete implements ArticleRepository.
func (a *AuthorRepo) DeleteById(ctx context.Context, id string) error {
	_, err := a.q.Author.WithContext(ctx).Where(a.q.Author.ID.Eq(id)).Delete()
	return err
}
