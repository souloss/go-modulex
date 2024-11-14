package routes

import (
	"net/http"

	"github.com/souloss/go-clean-arch/internal/usecase"
	"github.com/souloss/go-clean-arch/pkg/server"
)

func GetArticleRouteGroup(ArtilceUsecase *usecase.ArticleUsecase) *server.RouteGroup {
	rg := server.NewGroup("/artilct", server.WithRoutes(
		*server.NewRoute(http.MethodGet, "", ArtilceUsecase.GetAll),
		*server.NewRoute(http.MethodGet, "/:id", ArtilceUsecase.GetById),
		*server.NewRoute(http.MethodPost, "", ArtilceUsecase.Save),
		*server.NewRoute(http.MethodPut, "", ArtilceUsecase.UpdateById),
		*server.NewRoute(http.MethodDelete, "", ArtilceUsecase.DeleteById),
	))
	return rg
}
