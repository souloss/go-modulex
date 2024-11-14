package entities

import (
	"context"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

// Article is representing the Article data struct
type Article struct {
	ID        string    `json:"id" gorm:"type:char(36);primarykey" description:"文章ID"`
	Title     string    `json:"title" validate:"required" description:"文章标题"`
	Content   string    `json:"content" validate:"required" description:"文章内容"`
	Author    Author    `json:"author" gorm:"index;foreignKey:id" description:"文章作者"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at" description:"更新时间"`
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at" description:"创建时间"`
}

func (a *Article) TableName() string {
	return "article"
}

func (a *Article) BeforeCreate(*gorm.DB) error {
	a.ID = xid.New().String()
	return nil
}

type ArticleByIDRequest struct {
	ID string `json:"id"`
}

type UpdateArticleRequest struct {
	ArticleByIDRequest
	Title   string `json:"title"`
	Content string `json:"content" validate:"required"`
	Author  Author `json:"author"`
}

type SaveArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content" validate:"required"`
	Author  Author `json:"author"`
}

type ArticleRepository interface {
	GetAll(ctx context.Context) ([]*Article, error)
	GetById(ctx context.Context, id string) (*Article, error)
	Save(ctx context.Context, article *Article) error
	UpdateById(ctx context.Context, id string, r *UpdateArticleRequest) (*Article, error)
	DeleteById(ctx context.Context, id string) error
}

type ArticleCache interface {
	GetById(ctx context.Context, id string) (*Article, error)
}

type ArticleUsecase interface {
	GetAll(ctx context.Context) ([]*Article, error)
	GetById(ctx context.Context, req *ArticleByIDRequest) (*Article, error)
	Save(ctx context.Context, r *SaveArticleRequest) error
	UpdateById(ctx context.Context, r *UpdateArticleRequest) (*Article, error)
	DeleteById(ctx context.Context, req *ArticleByIDRequest) error
}
