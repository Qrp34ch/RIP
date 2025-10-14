package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "lab1/cmd/webserver/docs"
	"lab1/internal/app/config"
	"lab1/internal/app/repository"
)

type Handler struct {
	Repository *repository.Repository
	Config     *config.Config
}

func NewHandler(c *config.Config, r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
		Config:     c,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/reaction", h.GetReactions)
	router.GET("/reaction/:id", h.GetReaction)
	router.GET("/synthesis/:id", h.GetSynthesis)
	router.POST("/add-reaction-in-synthesis", h.AddReactionInSynthesis)
	router.POST("/delete/:id", h.RemoveSynthesis)

	//API
	//домен услуги (реакции)
	router.GET("/API/reaction", h.GetReactionsAPI)
	router.GET("/API/reaction/:id", h.GetReactionAPI)
	router.POST("/API/create-reaction", h.CreateReactionAPI)
	router.PUT("/API/reaction/:id", h.ChangeReactionAPI)
	router.DELETE("/API/reaction/:id", h.DeleteReactionAPI)
	router.POST("/API/reaction/:id/add-reaction-in-synthesis", h.AddReactionInSynthesisAPI)
	router.POST("/API/reaction/:id/image", h.UploadReactionImageAPI)

	//домен заявки (синтез)
	router.GET("/API/synthesis/icon", h.GetSynthesisIconAPI)
	router.GET("/API/synthesis", h.GetSynthesesAPI)
	router.GET("/API/synthesis/:id", h.GetSynthesisAPI)
	router.PUT("/API/synthesis/:id", h.UpdateSynthesisPurityAPI)
	router.PUT("/API/synthesis/:id/form", h.FormSynthesisAPI)
	router.PUT("/API/synthesis/:id/moderate", h.CompleteOrRejectSynthesisAPI)
	router.DELETE("/API/synthesis", h.DeleteSynthesisAPI)

	//домен м-м
	router.DELETE("/API/reaction-synthesis", h.RemoveReactionFromSynthesisAPI)
	router.PUT("/API/reaction-synthesis", h.UpdateReactionInSynthesisAPI)

	//домен пользователь
	router.POST("/API/users/register", h.RegisterUserAPI)
	router.GET("/API/users/profile", h.GetUserProfileAPI) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	router.POST("/API/users/login", h.LoginUserAPI)
	router.POST("/API/users/logout", h.LogoutUserAPI)
	router.PUT("/API/users/profile", h.UpdateUserAPI)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
