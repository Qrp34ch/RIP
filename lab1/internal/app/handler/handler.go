package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"lab1/internal/app/repository"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetReactions(ctx *gin.Context) {
	var reactions []repository.Reaction
	var err error

	searchQuery := ctx.Query("query") // получаем значение из поля поиска
	if searchQuery == "" {            // если поле поиска пусто, то просто получаем из репозитория все записи
		reactions, err = h.Repository.GetReactions()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		reactions, err = h.Repository.GetReactionsByTitle(searchQuery) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}

	synthesisReactions, err := h.Repository.GetReactionsInSynthesis()
	synthesisCount := len(synthesisReactions)
	if err != nil {
		logrus.Error(err)
		synthesisCount = 0
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":           time.Now().Format("15:04:05"),
		"reactions":      reactions,
		"query":          searchQuery, // передаем введенный запрос обратно на страницу
		"synthesisCount": synthesisCount,
		// в ином случае оно будет очищаться при нажатии на кнопку
	})
}

func (h *Handler) GetReaction(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	reaction, err := h.Repository.GetReaction(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "reaction.html", gin.H{
		"reaction": reaction,
	})
}

func (h *Handler) GetReactionsInSynthesis(ctx *gin.Context) {
	var reactions []repository.Reaction
	var err error

	reactions, err = h.Repository.GetReactionsInSynthesis()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "synthesis.html", gin.H{
		"reactions": reactions,
	})
}
