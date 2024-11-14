package routes

import (
	"net/http"

	"github.com/souloss/go-clean-arch/internal/usecase"
	"github.com/souloss/go-clean-arch/pkg/server"
)

func GetAuthorRouteGroup(AuthorUsecase *usecase.AuthorUsecase) *server.RouteGroup {
	rg := server.NewGroup("/author", server.WithRoutes(
		*server.NewRoute(http.MethodGet, "", AuthorUsecase.GetAll),
		*server.NewRoute(http.MethodGet, "/:id", AuthorUsecase.GetById),
		*server.NewRoute(http.MethodPost, "", AuthorUsecase.Save),
		*server.NewRoute(http.MethodPut, "", AuthorUsecase.UpdateById),
		*server.NewRoute(http.MethodDelete, "", AuthorUsecase.DeleteById),
	))
	return rg
}
