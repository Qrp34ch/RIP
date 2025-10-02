package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"lab1/internal/app/repository"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты, чтобы не писать все в одном месте
func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/reaction", h.GetReactions)
	router.GET("/reaction/:id", h.GetReaction)
	router.GET("/synthesis/:id", h.GetSynthesis)
	router.POST("/add-reaction-in-synthesis", h.AddReactionInSynthesis)
	router.POST("/delete/:id", h.RemoveSynthesis)

	//домен услуги (реакций)
	router.GET("/API/reaction", h.GetReactionsAPI)
	router.GET("/API/reaction/:id", h.GetReactionAPI)
	router.POST("/API/create-reaction", h.CreateReactionAPI)
	router.PUT("/API/reaction/:id", h.ChangeReactionAPI)
	router.DELETE("/API/reaction/:id", h.DeleteReactionAPI)
	router.POST("/API/reaction/:id/add-reaction-in-synthesis", h.AddReactionInSynthesisAPI)
	router.POST("/API/reaction/:id/image", h.UploadReactionImageAPI)

}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}

// errorHandler для более удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
