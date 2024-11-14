package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/souloss/go-clean-arch/pkg/server"
	"github.com/souloss/go-clean-arch/pkg/version"
)

type VersionResp struct {
	Version     string `json:"version"`
	BuildTime   string `json:"build_time"`
	GitBranch   string `json:"git_branch"`
	GitRevision string `json:"git_revision"`
}

func GetVersionRouter() server.Route {
	return server.Route{
		Method: http.MethodGet,
		Path:   "/version",
		Handler: func(c context.Context) (interface{}, error) {
			fmt.Println(VersionResp{
				Version:     version.Version,
				BuildTime:   version.BuildTime,
				GitBranch:   version.GitBranch,
				GitRevision: version.GitRevision,
			})
			return VersionResp{
				Version:     version.Version,
				BuildTime:   version.BuildTime,
				GitBranch:   version.GitBranch,
				GitRevision: version.GitRevision,
			}, nil
		},
	}
}
