package router

import (
	"github.com/GEtBUsyliVn/url-shortener/services/api-gateway/rest/handler"
	"github.com/gin-gonic/gin"
)

type Router struct {
	router  *gin.Engine
	handler *handler.Handler
}

func NewRouter(engine *gin.Engine, handler *handler.Handler) *Router {
	return &Router{
		router:  engine,
		handler: handler,
	}
}

func (r *Router) RegisterRoutes() {
	r.router.GET("/:code", r.handler.GetUrl)
	api := r.router.Group("/api")
	api.POST("/shorten", r.handler.CreateShortCode)
	api.GET("/stats/:code", r.handler.GetStats)

}
