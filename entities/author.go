package entities

import (
	"context"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

// Author representing the Author data struct
type Author struct {
	ID        string    `json:"id" gorm:"type:char(36);primarykey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (a *Author) TableName() string {
	return "author"
}

func (a *Author) BeforeCreate(*gorm.DB) error {
	a.ID = xid.New().String()
	return nil
}

type AuthoryIDRequest struct {
	ID string `json:"id"`
}

type AuthoryNameRequest struct {
	Name string `json:"name"`
}

type UpdateAuthorRequest struct {
	AuthoryIDRequest
	AuthoryNameRequest
}

type AuthorRepository interface {
	GetAll(ctx context.Context) ([]*Author, error)
	GetById(ctx context.Context, id string) (*Author, error)
	Save(ctx context.Context, name string) error
	UpdateById(ctx context.Context, id string, name string) (*Author, error)
	DeleteById(ctx context.Context, id string) error
}

type AuthorUsecase interface {
	GetAll(ctx context.Context) ([]*Author, error)
	GetById(ctx context.Context, req *AuthoryIDRequest) (*Author, error)
	Save(ctx context.Context, req *AuthoryNameRequest) error
	UpdateById(ctx context.Context, req *UpdateAuthorRequest) (*Author, error)
	DeleteById(ctx context.Context, req *AuthoryIDRequest) error
}
