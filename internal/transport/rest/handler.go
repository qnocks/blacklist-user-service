package rest

import (
	"github.com/gin-gonic/gin"
	_ "github.com/qnocks/blacklist-user-service/docs"
	"github.com/qnocks/blacklist-user-service/internal/service"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	h.initAPI(router)
	h.initSwagger(router)

	return router
}

func (h *Handler) initSwagger(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (h *Handler) initAPI(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", h.login)
	}

	blacklist := router.Group("/api/blacklist", h.authTokenMiddleware)
	{
		blacklist.POST("/", h.saveBlacklistedUser)
		blacklist.DELETE("/:id", h.deleteBlacklistedUser)
		blacklist.GET("/", h.getBlacklistedUsers)
	}
}
