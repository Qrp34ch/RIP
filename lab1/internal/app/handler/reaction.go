package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"lab1/internal/app/ds"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) GetReactions(ctx *gin.Context) {
	var reactions []ds.Reaction
	var err error

	searchQuery := ctx.Query("query") // получаем значение из нашего поля
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

	fmt.Printf("яхочупитсыномер2", h.Repository.FindUserSynthesis(1))

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"synthesisCount": h.Repository.GetReactionsInSynthesis(),
		"reactions":      reactions,
		"id":             h.Repository.FindUserSynthesis(1),
		"query":          searchQuery, // передаем введенный запрос обратно на страницу
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

func (h *Handler) AddReactionInSynthesis(ctx *gin.Context) {
	// считываем значение из формы, которую мы добавим в наш шаблон
	strId := ctx.PostForm("reaction_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	// Вызов функции добавления чата в заявку
	err = h.Repository.AddReactionInSynthesis(uint(id))
	if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return
	}

	// после вызова сразу произойдет обновление страницы
	ctx.Redirect(http.StatusFound, "/reaction")
}

func (h *Handler) GetSynthesis(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}
	synthesisStatus, err := h.Repository.SynthesisStatusById(uint(id))
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/reaction")
	}

	// если заявка по которой переходим удалена, то перенаправляем на главную
	if synthesisStatus == "удалён" {
		ctx.Redirect(http.StatusFound, "/reaction")
	}
	var reactions []ds.Reaction
	reactions, err = h.Repository.GetSynthesis(uint(id))
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "synthesis.html", gin.H{
		"reactions": reactions,
		"id":        id,
		"user":      h.Repository.GetUserNameByID(1),
		"date":      h.Repository.GetDateUpdate(uint(id)),
	})
}

func (h *Handler) RemoveSynthesis(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	err = h.Repository.RemoveSynthesis(uint(id))
	ctx.Redirect(http.StatusFound, "/reaction")
}
