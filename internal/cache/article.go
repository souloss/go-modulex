package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/souloss/go-clean-arch/entities"
)

type articleCache struct {
	repo  entities.ArticleRepository
	cache *redis.Client
}

// 确保实现了接口
var _ entities.ArticleCache = (*articleCache)(nil)

// GetById implements entities.ArticleCache.
func (a *articleCache) GetById(ctx context.Context, id string) (*entities.Article, error) {
	var article *entities.Article
	if err := a.cache.HGetAll(ctx, id).Scan(&article); err == nil {
		return article, nil
	}

	article, err := a.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := a.cache.HSet(ctx, id, article).Err(); err != nil {
		fmt.Println("error")
	}
	return article, nil
}
