package main

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/souloss/go-clean-arch/pkg/server"
	"github.com/souloss/go-clean-arch/pkg/server/gin"
)

type RequestData struct {
	ID      int    `uri:"id"`                           // 来自Path
	Name    string `header:"X-Name" binding:"required"` // 来自Header
	Age     int    `form:"age"`                         // 来自Form Data (通常是POST请求的Body)
	Email   string `form:"email"`                       // 来自Query String
	Message string `json:"message"`                     // 来自JSON Body
}

func MYController(ctx context.Context, req *RequestData) (string, error) {
	spew.Dump(req)
	return "hello", nil
}

func main() {
	g := gin.NewGinFramework()
	g.Init(&server.ServerConfig{
		Port: 8080,
	}, &server.Options{})
	g.RegisterRoutes([]server.Route{
		{
			Method:      "GET",
			Path:        "/hello/:id",
			Handler:     MYController,
			Middlewares: []interface{}{},
		},
		{
			Method:      "POST",
			Path:        "/hello/:id",
			Handler:     MYController,
			Middlewares: []interface{}{},
		},
	})
	g.Start(context.TODO())

}

// curl -XPOST -d '{"message": "Hello, Worl","age":123}'  -H "X-Name: John Doe" http://localhost:8080/path/query/bind/5?email=123 -v
